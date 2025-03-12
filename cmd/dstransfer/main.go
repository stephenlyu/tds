package main

import (
	"bufio"
	"flag"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stephenlyu/tds/datasource"
	csvdatasource "github.com/stephenlyu/tds/datasource/csv"
	mongodatasource "github.com/stephenlyu/tds/datasource/mongo"
	redisdatasource "github.com/stephenlyu/tds/datasource/redis"
	tdxdatasource "github.com/stephenlyu/tds/datasource/tdx"
	"github.com/stephenlyu/tds/date"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/tradedate"
	"github.com/stephenlyu/tds/util"
)

//
// 将数据在数据源间移动
//

func moveData(srcDs datasource.BaseDataSource, destDs datasource.BaseDataSource,
	security *entity.Security, p period.Period, startDate string, endDate string,
	useQfq bool) error {
	var start, end uint64
	if startDate != "" {
		startTs, _, _, _ := tradedate.GetTradeDateRangeByDateString(security, startDate)
		start, _ = date.SecondString2Timestamp(startTs)
	}

	if endDate != "" {
		_, endTs, _, _ := tradedate.GetTradeDateRangeByDateString(security, endDate)
		end, _ = date.SecondString2Timestamp(endTs)
	}

	var err error
	var data []entity.Record

	if useQfq {
		err, data = srcDs.GetForwardAdjustedRangeData(security, p, start, end)
	} else {
		err, data = srcDs.GetRangeData(security, p, start, end)
	}
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
	case "tdx":
		tdxDir := "data"
		if len(params) > 0 {
			tdxDir = params[0]
		}
		return tdxdatasource.NewDataSource(tdxDir, true)
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

func readLines(filename string) ([]string, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// 确保文件在函数结束时关闭
	defer file.Close()

	var lines []string
	// 创建一个新的扫描器，用于逐行读取文件
	scanner := bufio.NewScanner(file)
	// 逐行扫描文件
	for scanner.Scan() {
		// 将当前行添加到切片中
		lines = append(lines, scanner.Text())
	}
	// 检查扫描过程中是否发生错误
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func main() {
	srcParamPtr := flag.String("src", "csv|csv", "Source data source parameter")
	destParamPtr := flag.String("dest", "redis|localhost:6379|", "Destination data source parameter")
	codesPtr := flag.String("code", "", "Security codes")
	codesFilePtr := flag.String("code-file", "", "Security codes file")
	periodPtr := flag.String("period", "M1", "Data period")
	startDatePtr := flag.String("start-date", "", "Start date")
	endDatePtr := flag.String("end-date", "", "End date")
	useQfqPtr := flag.Bool("qfq", false, "Use QianFuQaun?")
	flag.Parse()

	err, p := period.PeriodFromString(*periodPtr)
	if err != nil {
		logrus.Fatal(err)
	}

	// Parse securities
	var codes []string
	if *codesFilePtr != "" {
		codes, err = readLines(*codesFilePtr)
		if err != nil {
			logrus.Fatalf("read code file fail, err:%+v", err)
		}
	} else {
		codes = strings.Split(*codesPtr, ",")
	}
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
		err = moveData(srcDs, destDs, security, p, *startDatePtr, *endDatePtr, *useQfqPtr)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}
