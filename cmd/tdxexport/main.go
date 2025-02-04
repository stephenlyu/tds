package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	tdxdatasource "github.com/stephenlyu/tds/datasource/tdx"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/period"
	"github.com/z-ray/log"
)

type R struct {
	Date   string  `json:"date"`
	Open   float32 `json:"open"`
	Close  float32 `json:"close"`
	High   float32 `json:"high"`
	Low    float32 `json:"low"`
	Volume float32 `json:"volume"`
	Amount float32 `json:"amount"`
}

func main() {
	dataDir := flag.String("data-dir", "", "通达信数据目录")
	periodStr := flag.String("period", "d1", "周期")
	code := flag.String("code", "", "股票代码")

	flag.Parse()

	ds := tdxdatasource.NewDataSource(*dataDir, true)
	security, err := entity.ParseSecurity(*code)
	if err != nil {
		log.Fatalf("错误：%s", err.Error())
	}

	_, period := period.PeriodFromString(*periodStr)
	err, data := ds.GetData(security, period)
	if err != nil {
		log.Fatalf("加载数据失败，错误：%s", err.Error())
	}

	result := make([]R, len(data))
	for i := range data {
		r := &data[i]
		rr := &result[i]
		rr.Date = r.GetDate()
		rr.Open = float32(r.GetOpen())
		rr.Close = float32(r.GetClose())
		rr.High = float32(r.GetHigh())
		rr.Low = float32(r.GetLow())
		rr.Volume = float32(r.GetVolume())
		rr.Amount = float32(r.GetAmount())
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	ioutil.WriteFile(fmt.Sprintf("%s.json", *code), bytes, 0666)
}
