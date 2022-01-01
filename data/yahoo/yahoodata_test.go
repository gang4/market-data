package yahoo

import (
	"fmt"
	"market-data/algorithms"
	"testing"
	"time"
)

func TestGetData(t *testing.T) {
	data, err := GetYahooData("iwm", "2d", "1d")
	meta := data.Chart.Result[0].Meta
	quote := data.Chart.Result[0].Indicators.Quote
	adjust := data.Chart.Result[0].Indicators.Adjclose

	fmt.Println(meta)
	fmt.Println(quote)
	fmt.Println(adjust)
	fmt.Println(err)
}

func TestGetData1(t *testing.T) {
	period := "15d"
	data, err := GetYahooData("iwm", period, "1d")
	if err != nil {
		panic(err)
	}
	max, index := algorithms.MaxProfit(data.Chart.Result[0].Indicators.Adjclose[0].Adjclose)
	fmt.Println("Max: ", max)
	date := time.Unix(int64(data.Chart.Result[0].Timestamp[index]), 0)
	fmt.Println("Date: ", date)
}
