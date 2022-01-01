package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"market-data/cmd/config"
	"strconv"
	"strings"
	"sync"
	"time"

	"market-data/data/yahoo"

	mysqlErrors "github.com/go-mysql/errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

const (
	username = "root"
	password = "9005"
	hostname = "127.0.0.1:3306"
	dbname   = "mydb"
)

// Query query from db
type Query interface {
	query(db *sql.DB, query string) error
}

type QueryType int

const (
	Select = iota
	Create
	Insert
	ShowTable
	LatestRecord
)

// Resp is a table.
type Entity struct {
	Open    float64
	Close   float64
	High    float64
	Low     float64
	Adjust  float64
	Vol     int64
	TimeStr string
}

// Req query
type Req struct {
	ReqType QueryType
	Query   string
	Params  []Entity
}

type Resp interface{}

// QueryRequest communicated from clients to here
type QueryRequest interface {
	OnNext(entities Resp)
	OnComplete()
	OnError(err error)
	GetReq() *Req
}

// QuerySQLDB struct for mutation
type QuerySQLDB struct {
	request  chan QueryRequest
	finished chan bool
	cfg      *config.Config
}

// DB export for clients to use
var DB QuerySQLDB

func (d *QuerySQLDB) dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

// Initialize  -- setup
func Initialize(cfg *config.Config) {
	DB = QuerySQLDB{
		make(chan QueryRequest, cfg.ReqQueueSize),
		make(chan bool, 1),
		cfg}
	go DB.queryLoop()
}

func (d *QuerySQLDB) queryLoop() error {
	// "mysql" is driver name
	db, err := sql.Open("mysql", d.dsn())
	if err != nil {
		return err
	}
	defer db.Close()
	db.SetMaxOpenConns(d.cfg.MaxOpenConns)
	db.SetMaxIdleConns(d.cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB: ", err)
		return err
	}
	glog.Infof("Connected to DB %s successfully\n", dbname)

	d.query(db)

	return nil
}

func (d *QuerySQLDB) query(db *sql.DB) {
	for {
		select {
		case <-d.finished:
			// fmt.Println("exiting query")
			return
		case req := <-d.request:
			// fmt.Println("request: ", req)
			d.handleReq(db, req)
		case <-time.After(time.Second * 60):
			// Our metrics are 10s
			glog.Infoln("No writes for 1 minute")
			continue
		}
	}
}

func (d *QuerySQLDB) handleReq(db *sql.DB, req QueryRequest) error {
	defer req.OnComplete()

	switch req.GetReq().ReqType {
	case Select:
		resp, err := d.handleSelect(db, req)
		if err == nil {
			req.OnNext(resp)
		} else {
			req.OnError(err)
		}
	case Create:
		err := d.handleCreate(db, req)
		if err == nil {
			req.OnNext(nil)
		} else {
			req.OnError(err)
		}
	case Insert:
		err := d.handleInsert(db, req)
		if err == nil {
			req.OnNext(nil)
		} else {
			req.OnError(err)
		}
	case ShowTable:
		name, err := d.handleShowTable(db, req)
		if err == nil {
			req.OnNext(name)
		} else {
			req.OnError(err)
		}
	case LatestRecord:
		latestRecord, err := d.handleGetLatestRecord(db, req)
		if err == nil {
			req.OnNext(latestRecord)
		} else {
			req.OnError(err)
		}
	default:
		return errors.New("Bad Query Type")
	}
	return nil
}

func (d *QuerySQLDB) handleSelect(db *sql.DB, reqInt QueryRequest) ([]Entity, error) {
	// fmt.Println("handleSelect")
	req := reqInt.GetReq()

	selDB, err := db.Query(req.Query)
	if err != nil {
		return nil, err
	}
	defer selDB.Close()
	var open, close, high, low, adjust float64
	var vol int64
	var timeStr string
	var rts []Entity = make([]Entity, 0)
	for selDB.Next() {
		err = selDB.Scan(&timeStr, &open, &close, &high, &low, &adjust, &vol)
		if err != nil {
			return nil, err
		} else {
			rt := Entity{}
			rt.Open = open
			rt.Close = close
			rt.High = high
			rt.Low = low
			rt.Adjust = adjust
			rt.Vol = vol
			rt.TimeStr = timeStr
			rts = append(rts, rt)
		}
	}
	return rts, nil
}

func (d *QuerySQLDB) handleCreate(db *sql.DB, reqInt QueryRequest) error {
	// fmt.Println("handleCreate")
	sql := "CREATE TABLE IF NOT EXISTS " + reqInt.GetReq().Query + "(" +
		"`timestamp` VARCHAR(10) NOT NULL, " +
		"`Open` DOUBLE NOT NULL, " +
		"`close` DOUBLE NOT NULL, " +
		"`high` DOUBLE NOT NULL, " +
		"`low` DOUBLE NOT NULL, " +
		"`adjust` DOUBLE NULL, " +
		"`volume` BIGINT NULL, " +
		"PRIMARY KEY (`timestamp`), " +
		"UNIQUE INDEX `timestamp_UNIQUE` (`timestamp` ASC) VISIBLE) " +
		"ENGINE = InnoDB"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (d *QuerySQLDB) handleShowTable(db *sql.DB, reqInt QueryRequest) (string, error) {
	// fmt.Println("handleShowTable")
	req := reqInt.GetReq()
	selDB, err := db.Query(req.Query)
	if err != nil {
		return "", err
	}
	defer selDB.Close()
	var name string
	for selDB.Next() {
		err = selDB.Scan(&name)
		if err != nil {
			return "", err
		}
	}
	return name, nil
}

func (d *QuerySQLDB) handleInsert(db *sql.DB, reqInt QueryRequest) error {
	// fmt.Println("handleInsert")
	req := reqInt.GetReq()

	stmt, err := db.Prepare(req.Query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := range req.Params {
		_, err := stmt.Exec(req.Params[i].TimeStr, req.Params[i].Open, req.Params[i].Close, req.Params[i].High, req.Params[i].Low, req.Params[i].Adjust, req.Params[i].Vol)
		if err != nil {
			if mysqlErrors.MySQLErrorCode(err) == 1062 {
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

func (d *QuerySQLDB) handleGetLatestRecord(db *sql.DB, reqInt QueryRequest) (string, error) {
	// fmt.Println("handleGetLatestRecord")
	selDB, err := db.Query(reqInt.GetReq().Query)
	if err != nil {
		return "", err
	}
	defer selDB.Close()
	var timestamp string
	for selDB.Next() {
		err = selDB.Scan(&timestamp)
		if err != nil {
			return "", err
		}
		break
	}
	return timestamp, nil
}

// SendReq -- client to send request
func (d *QuerySQLDB) SendReq(req QueryRequest) error {
	select {
	case d.request <- req:
	default:
		err := errors.New("buffer full while calling API")
		glog.Infoln(err)
		return err
	}

	return nil
}

// Stop -- client to send request
func (d *QuerySQLDB) Stop() {
	close(d.finished)
}

type queryResp interface{}

type client struct {
	req          *Req
	wg           *sync.WaitGroup
	latestRecord queryResp
	err          error
}

func (c *client) OnNext(resp Resp) {
	if c.req.ReqType == LatestRecord {
		en, ok := resp.(string)
		if !ok {
			c.err = errors.New("Invalid return for last record")
		} else {
			c.latestRecord = en
		}
	} else if c.req.ReqType == Select {
		c.latestRecord = resp
	}
}

func (c *client) OnError(err error) {
	c.err = err
	// fmt.Println(err)
}

func (c *client) GetReq() *Req {
	return c.req
}

func (c *client) OnComplete() {
	c.wg.Done()
}

func (d *QuerySQLDB) clientQuery(req *Req) (interface{}, error) {
	var wg sync.WaitGroup
	wg.Add(1)
	var c client = client{req, &wg, "", nil}
	DB.SendReq(&c)
	wg.Wait()
	if c.err != nil {
		return nil, c.err
	}
	return c.latestRecord, nil
}

func (d *QuerySQLDB) tableExist(name string) (bool, error) {
	req := Req{}
	req.ReqType = LatestRecord
	req.Query = "select timestamp from " + "mydb." + name + " order by timestamp DESC limit 1"
	_, err := d.clientQuery(&req)
	if err != nil {
		if mysqlErrors.MySQLErrorCode(err) == 1146 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *QuerySQLDB) getTableNumberRows(name string) (int, error) {
	req := Req{}
	req.ReqType = LatestRecord
	req.Query = "SELECT COUNT(*) FROM mydb." + name
	row, err := d.clientQuery(&req)
	if err != nil {
		return 0, err
	}
	rows, err1 := strconv.ParseInt(row.(string), 10, 32)
	if err1 != nil {
		return 0, err1
	}
	return int(rows), nil
}

func (d *QuerySQLDB) getLastRecord(symbol string) (int, error) {
	req := Req{}
	req.ReqType = LatestRecord
	req.Query = "select timestamp from " + "mydb." + symbol + " order by timestamp DESC limit 1"
	rawString, err := d.clientQuery(&req)
	if err != nil {
		return 0, err
	}
	if len(rawString.(string)) == 0 {
		// Table is empty, load at least one row
		return 1, nil
	}
	// Is db current
	days, err := DaysFromToday(rawString.(string))
	if err != nil {
		return 0, err
	}
	return days, nil
}

// EnsureTableExist --
func (d *QuerySQLDB) EnsureTableExist(symbol string) error {
	exist, err := d.tableExist(symbol)
	if err != nil {
		return err
	}
	if !exist {
		// create table
		req := Req{}
		req.ReqType = Create
		req.Query = "mydb." + symbol
		var wg sync.WaitGroup
		wg.Add(1)
		_, err = d.clientQuery(&req)
		if err != nil {
			return nil
		}
	}
	return nil
}

// LoadDataToDB -- load record to db if not update to today
func (d *QuerySQLDB) LoadDataToDB(symbol string, dayRange string, interval string) error {
	r := strings.Trim(dayRange, "d")
	required, err := strconv.ParseInt(r, 10, 32)
	if err != nil {
		return err
	}
	rows, err1 := d.getTableNumberRows(symbol)
	if err1 != nil {
		return err1
	}

	days, err2 := d.getLastRecord(symbol)
	if err2 != nil {
		return err2
	}
	if rows < int(required) {
		if days < int(required) {
			days = int(required)
		}
	}
	t := time.Now()
	if t.Weekday() == time.Sunday {
		days -= 2
	} else if t.Weekday() == time.Saturday {
		days--
	}
	if days == 0 {
		// assume there is no gap in table
		return nil
	}

	// always fill the gap
	i := strconv.FormatInt(int64(days), 10)
	quotes, err3 := yahoo.GetYahooData(symbol, i+"d", interval)
	if err3 != nil {
		return err3
	}
	if quotes.Chart.Result == nil {
		// Holiday?
		return nil
	}
	// Insert difference into db
	req := Req{}
	req.ReqType = Insert
	req.Query = "INSERT INTO mydb." + symbol + " (timestamp, Open, close, high, low, adjust, volume) values (?, ?, ?, ?, ?, ?, ?)"
	fillParams(quotes, &req.Params)
	_, err = d.clientQuery(&req)
	if err != nil {
		return err
	}
	return nil
}

// LoadData -- load required rows
func (d *QuerySQLDB) LoadData(symbol string, dayRange string, interval string) ([]Entity, error) {
	err := d.EnsureTableExist(symbol)
	if err != nil {
		return nil, err
	}
	err = d.LoadDataToDB(symbol, dayRange, interval)
	if err != nil {
		return nil, err
	}
	r := strings.Trim(dayRange, "d")
	req := Req{}
	req.ReqType = Select
	req.Query = "select * from mydb." + symbol + " order by timestamp DESC limit " + r
	rawResp, err1 := d.clientQuery(&req)
	if err1 != nil {
		return nil, err1
	}
	return rawResp.([]Entity), nil
}

func (d *QuerySQLDB) LoadSimpleMovingAverge(symbol string, dayRange string, period string, interval string) ([]float64, error) {
	dayInt, err := strconv.Atoi(dayRange)
	pp := strings.Trim(period, "d")
	length, err1 := strconv.Atoi(pp)
	if err == nil && err1 == nil {
		l := dayInt
		if length < dayInt {
			length = dayInt
		}
		dayInt = dayInt + length
		p := strconv.Itoa(dayInt)
		data, err := d.LoadData(symbol, p+"d", interval)
		// for i := 0; i < len(data); i++ {
		// 	fmt.Println(data[i].Close)
		// }
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		pData := make([]float64, length)
		var base float64 = 0.0
		for i := 0; i < (l - 1); i++ {
			base += data[i].Close
		}

		for i := 0; i < length; i++ {
			pData[i] = base + data[l-1+i].Adjust
			base = pData[i] - data[i].Adjust
			pData[i] = pData[i] / float64(l)
		}
		// fmt.Println(pData)
		return pData, nil
	}
	return nil, err
}

func fillParams(data *yahoo.AutoGenerated, params *[]Entity) {
	timestamp := data.Chart.Result[0].Timestamp
	quote := data.Chart.Result[0].Indicators.Quote
	adjust := data.Chart.Result[0].Indicators.Adjclose
	*params = make([]Entity, len(quote[0].Open))
	for i := range quote[0].Open {
		(*params)[i].Open = quote[0].Open[i]
		(*params)[i].Close = quote[0].Close[i]
		(*params)[i].High = quote[0].High[i]
		(*params)[i].Low = quote[0].Low[i]
		(*params)[i].Vol = int64(quote[0].Volume[i])
		(*params)[i].Adjust = adjust[0].Adjclose[i]
		(*params)[i].TimeStr = getTimeString(timestamp[i])
	}
}

func getTimeString(tt int) string {
	date := time.Unix(int64(tt), 0)
	yy := date.Year()
	mm := date.Month()
	dd := date.Day()
	s := fmt.Sprintf("%d-%02d-%02d", yy, mm, dd)
	return s
}
