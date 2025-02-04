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
)

//
// 将数据在数据源间移动
//

func moveData(srcDs datasource.BaseDataSource, destDs datasource.BaseDataSource, security *entity.Security, p period.Period, startDate string, endDate string) error {
	var start, end uint64
	if startDate != "" {
		startTs, _, _, _ := tradedate.GetTradeDateRangeByDateString(security, startDate)
		start, _ = date.SecondString2Timestamp(startTs)
	}

	if endDate != "" {
		_, endTs, _, _ := tradedate.GetTradeDateRangeByDateString(security, endDate)
		end, _ = date.SecondString2Timestamp(endTs)
	}

	err, data := srcDs.GetRangeData(security, p, start, end)
	if err != nil {
		return err
	}

	err = destDs.SaveData(security, p, data)
	if err == nil {
		logrus.Infof("%d records transferred.", len(data))
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
		if len(params) > 2 {
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
	srcParamPtr := flag.String("src-param", "csv|csv", "Source data source parameter")
	destParamPtr := flag.String("dest-param", "redis|localhost:6379|", "Destination data source parameter")
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

	srcParams := strings.Split(*srcParamPtr, "|")
	destParams := strings.Split(*destParamPtr, "|")

	if len(srcParams) == 0 {
		logrus.Fatalf("Bad src-param %s", *srcParamPtr)
	}

	if len(destParams) == 0 {
		logrus.Fatalf("Bad dest-param %s", *destParamPtr)
	}

	srcDs := createDataSource(srcParams[0], srcParams[1:])
	destDs := createDataSource(destParams[0], destParams[1:])

	for i := range securities {
		security := &securities[i]
		err = moveData(srcDs, destDs, security, p, *startDatePtr, *endDatePtr)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}
