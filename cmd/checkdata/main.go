package main

import (
	"github.com/stephenlyu/tds/datasource"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/date"
	"github.com/stephenlyu/tds/entity"
	"flag"
	"github.com/stephenlyu/tds/util"
	"github.com/sirupsen/logrus"
	"strings"
	"github.com/stephenlyu/tds/datasource/csv"
	"github.com/stephenlyu/tds/datasource/redis"
	"github.com/stephenlyu/tds/datasource/mongo"
	"github.com/stephenlyu/tds/tradedate"
	"github.com/deckarep/golang-set"
	"sort"
)

//
// 将数据在数据源间移动
//

func checkData(ds datasource.BaseDataSource, security *entity.Security, p period.Period, startDate string, endDate string) error {
	var start, end uint64
	var startTs, endTs string
	if startDate != "" {
		startTs, _, _, _ = tradedate.GetTradeDateRangeByDateString(security, startDate)
		start, _ = date.SecondString2Timestamp(startTs)
	}

	if endDate != "" {
		_, endTs, _, _ = tradedate.GetTradeDateRangeByDateString(security, endDate)
		end, _ = date.SecondString2Timestamp(endTs)
	}

	err, data := ds.GetRangeData(security, p, start, end - 1)
	if err != nil {
		return err
	}

	dateRecords := make(map[string][]*entity.Record)

	var tradeDates []string
	var tradeDate string
	startTs, endTs, tradeDate, _ = tradedate.GetTradeDateRangeByDateString(security, data[0].GetDate())
	tradeDates = append(tradeDates, tradeDate)
	for i := range data {
		d := data[i].GetDate()
		if d >= endTs {
			startTs, endTs, tradeDate, _ = tradedate.GetTradeDateRangeByDateString(security, d)
			tradeDates = append(tradeDates, tradeDate)
		}
		dateRecords[tradeDate] = append(dateRecords[tradeDate], &data[i])
	}

	for _, d := range tradeDates {
		records := dateRecords[d]
		ts, _ := date.DayString2Timestamp(d)
		tickers := tradedate.GetTradeTickers(security, ts)
		tickerSet := mapset.NewSet()
		for _, ticker := range tickers {
			tickerSet.Add(ticker)
		}

		dateSet := mapset.NewSet()
		for i := range records {
			dateSet.Add(records[i].Date)
		}

		diff := tickerSet.Difference(dateSet)
		diffDates := diff.ToSlice()
		sort.SliceStable(diffDates, func(i, j int) bool {
			return diffDates[i].(uint64) < diffDates[j].(uint64)
		})

		for _, ts := range diffDates {
			logrus.Infof("%s", date.Timestamp2SecondString(ts.(uint64)))
		}
	}

	return err
}

func createDataSource(dsType string, params []string) datasource.BaseDataSource {
	switch dsType {
	case "csv":
		csvDir := "csv"
		if len(params) > 0 {
			csvDir = params[0]
		}
		return csvdatasource.NewCSVDataSource(csvDir)
	case "redis":
		redisUrl := "localhost:6379"
		redisPass := ""
		if len(params) >= 2 {
			logrus.Fatalf("Bad redis params %s", strings.Join(params, "|"))
		}
		switch len(params) {
		case 1:
			redisUrl = params[0]
		case 2:
			redisUrl = params[0]
			redisPass = params[1]
		}
		return redisdatasource.NewRedisDataSource(redisUrl, redisPass)

	case "mongo":
		mongoUrl := "localhost/data"
		dbName := "data"
		if len(params) > 2 {
			logrus.Fatalf("Bad mongo params %s", strings.Join(params, "|"))
		}
		switch len(params) {
		case 1:
			mongoUrl = params[0]
		case 2:
			mongoUrl = params[0]
			dbName = params[1]
		}

		return mongodatasource.NewMongoDataSource(mongoUrl, dbName)
	}

	util.UnreachableCode()
	return nil
}

func main() {
	dsPtr := flag.String("ds", "csv|csv", "Data source descriptor, support csv|csv,redis|localhost:6379| and mongo|dev:pwd@localhost/db|db")
	codesPtr := flag.String("code", "", "Security codes")
	periodPtr := flag.String("period", "M1", "Data period")
	startDatePtr := flag.String("start-date", "", "Start date")
	endDatePtr := flag.String("end-date", "", "End date")
	flag.Parse()

	err, p := period.PeriodFromString(*periodPtr)
	if err != nil {
		logrus.Fatal(err)
	}

	// Parse securities
	codes := strings.Split(*codesPtr, ",")
	var securities []entity.Security
	for _, code := range codes {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}

		security, err := entity.ParseSecurity(code)
		if err != nil {
			logrus.Fatalf("Bad code %s", code)
		}
		securities = append(securities, *security)
	}
	if len(securities) == 0 {
		logrus.Fatal("code required")
	}

	dsParts := strings.Split(*dsPtr, "|")

	if len(dsParts) == 0 {
		logrus.Fatalf("Bad ds %s", *dsPtr)
	}

	ds := createDataSource(dsParts[0], dsParts[1:])

	for i := range securities {
		security := &securities[i]
		err = checkData(ds, security, p, *startDatePtr, *endDatePtr)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}
