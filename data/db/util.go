package db

import (
	"strconv"
	"strings"
	"time"
)

func GetToday() string {
	t := time.Now()
	rt := t.Format(time.RFC3339)
	date := strings.Split(rt, "T")
	return date[0]
}

func toDate(d string) *time.Time {
	//fmt.Println("toDate: ", d)
	parts := strings.Split(d, "-")
	yy, err1 := strconv.Atoi(parts[0])
	mm, err2 := strconv.Atoi(parts[1])
	dd, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return nil
	}
	loc, _ := time.LoadLocation("Local")
	src := time.Date(yy, time.Month(mm), dd, 0, 0, 0, 0, loc)
	return &src
}

func DaysFromToday(d string) (int, error) {
	src := toDate(d)
	dr := time.Now().Sub(*src)
	hrs := int(dr.Hours())
	days := int(hrs / 24)
	return days, nil
}
