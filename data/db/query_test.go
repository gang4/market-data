package db

import (
	"fmt"
	"market-data/cmd/config"
	"sync"
	"testing"

	"github.com/golang/glog"
)

const (
	configPath string = "../../cmd/config/config.toml"
)

type testReq struct {
	req *Req
	wg  *sync.WaitGroup
}

func (c *testReq) OnNext(resp Resp) {
	if c.req.ReqType == Select {
		en, ok := resp.([]Entity)
		if !ok {
			fmt.Println("Wrong type!")
		} else {
			fmt.Println(en)
		}
	} else if c.req.ReqType == ShowTable || c.req.ReqType == LatestRecord {
		en, ok := resp.(string)
		if !ok {
			fmt.Println("Wrong type!")
		} else {
			fmt.Println(en)
		}
	}
}

func (c *testReq) OnError(err error) {
	fmt.Println(err)
}

func (c *testReq) GetReq() *Req {
	return c.req
}

func (c *testReq) OnComplete() {
	c.wg.Done()
}

func TestSetup(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	r := Req{}
	r.ReqType = Select
	var tt testReq = testReq{&r, nil}
	DB.SendReq(&tt)
	r.ReqType = Insert
	DB.SendReq(&tt)
	r.ReqType = Create
	DB.SendReq(&tt)
	DB.Stop()
	glog.Flush()
}

func TestSelect(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	req := Req{}
	req.ReqType = Select
	req.Query = "select * from mydb.spx where timestamp like '%2020-10%'"
	var wg sync.WaitGroup
	wg.Add(1)
	var tr testReq = testReq{&req, &wg}
	DB.SendReq(&tr)

	wg.Wait()
	DB.Stop()
	glog.Flush()
}

func TestInsert(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	req := Req{}
	req.ReqType = Insert
	req.Query = "INSERT INTO mydb.spx (timestamp, Open, close, high, low, adjust, volume) values (?, ?, ?, ?, ?, ?, ?)"
	req.Params = make([]Entity, 2)
	req.Params[0].TimeStr = "2021-01-01"
	req.Params[1].TimeStr = "2021-10-01"
	var wg sync.WaitGroup
	wg.Add(1)
	var tr testReq = testReq{&req, &wg}
	DB.SendReq(&tr)
	wg.Wait()
	DB.Stop()
	glog.Flush()
}

// show tables from mydb like "spx"
func TestShowTable(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	req := Req{}
	req.ReqType = ShowTable
	req.Query = "show tables from mydb like 'spx'"
	var wg sync.WaitGroup
	wg.Add(1)
	var tr testReq = testReq{&req, &wg}
	DB.SendReq(&tr)
	wg.Wait()
	DB.Stop()
	glog.Flush()
}

func TestCreateTable(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	createTable("mydb.snow")
	DB.Stop()
	glog.Flush()
}

func createTable(tableName string) {
	req := Req{}
	req.ReqType = Create
	req.Query = tableName
	var wg sync.WaitGroup
	wg.Add(1)
	var tr testReq = testReq{&req, &wg}
	DB.SendReq(&tr)
	wg.Wait()
}

func TestLatestRecord(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	req := Req{}
	req.ReqType = LatestRecord
	req.Query = "select timestamp from mydb.spx order by timestamp DESC limit 1"
	var wg sync.WaitGroup
	wg.Add(1)
	var tr testReq = testReq{&req, &wg}
	DB.SendReq(&tr)
	wg.Wait()
	DB.Stop()
	glog.Flush()
}

func TestLastRecord(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	rows, err := DB.getLastRecord("slb")
	fmt.Println(rows, err)
	DB.Stop()
}

func TestLoadDataToDB(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	createTable("mydb.slb")

	err := DB.LoadDataToDB("slb", "50d", "1d")
	fmt.Println(err)
	DB.Stop()
}

func TestTableExist(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	tabelname := "abc"
	exists, err := DB.tableExist("abc")
	fmt.Println(tabelname, exists, err)
	DB.Stop()
}

func TestRows(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	rows, err := DB.getTableNumberRows("slb")
	fmt.Println(rows, err)
	DB.Stop()
}

func TestEnsureTableExist(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	err := DB.EnsureTableExist("pton")
	fmt.Println(err)
	DB.Stop()
}

func TestLoadData(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	data, err := DB.LoadData("ipob", "20d", "1d")
	DB.Stop()
	fmt.Println(data, err)
}

func TestLoadSimpleAverage(t *testing.T) {
	config.LoadConfigFromFile(configPath)
	Initialize(&config.Cfg)

	// e, err := DB.LoadData("arkf", "103d", "1d")
	// fmt.Println(len(e), err)

	data, err := DB.LoadSimpleMovingAverge("arkf", "3", "100d", "1d")
	fmt.Println(len(data), err)
}
