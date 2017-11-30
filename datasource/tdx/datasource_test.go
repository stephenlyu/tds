package tdxdatasource_test

import (
	"testing"
	"github.com/stephenlyu/tds/datasource/tdx"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tds/period"
	"github.com/z-ray/log"
	"github.com/stephenlyu/tds/date"
	"fmt"
)

func TestTdxDataSource(t *testing.T) {
	log.SetOutputLevel(log.Ldebug)

	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	err, period2 := period.PeriodFromString("M5")
	util.Assert(err == nil, "")

	err, period3 := period.PeriodFromString("D1")
	util.Assert(err == nil, "")

	err, period4 := period.PeriodFromString("M15")
	util.Assert(err == nil, "")

	startDate, err := date.DayString2Timestamp("20170101")
	util.Assert(err == nil, "")

	endDate, err := date.DayString2Timestamp("20170210")
	util.Assert(err == nil, "")

	ds := tdxdatasource.NewDataSource("data", true)

	err, items := ds.GetStockInfoEx(security)
	util.Assert(err == nil, "")
	util.Assert(len(items) == 21, "")

	// 分钟数据
	err, data := ds.GetData(security, period1)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 49920, "")

	err, data = ds.GetRangeData(security, period1, startDate, endDate)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 1200, "")

	// 5分钟数据
	err, data = ds.GetData(security, period2)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 22608, "")

	err, data = ds.GetRangeData(security, period2, startDate, endDate)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 1104, "")

	// 日线数据
	err, data = ds.GetData(security, period3)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 1571, fmt.Sprintf("got %d", len(data)))

	err, data = ds.GetRangeData(security, period3, startDate, endDate)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 24, "")

	// 15分钟数据
	err, data = ds.GetData(security, period4)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 7536, "")

	err, data = ds.GetRangeData(security, period4, startDate, endDate)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 368, "")

	// 日线数据前复权
	//startDate1, err := date.DayString2Timestamp("20170720")
	//util.Assert(err == nil, "")
	//
	//endDate1, err := date.DayString2Timestamp("20170721")
	//util.Assert(err == nil, "")

	err, data = ds.GetForwardAdjustedRangeData(security, period3, 0, 0)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 1571, "")
}

func TestTdxDataSource_GetStockCodes(t *testing.T) {
	ds := tdxdatasource.NewDataSource("data", true)
	stockCodes := ds.GetStockCodes("sz")
	fmt.Println(len(stockCodes))
	//for _, c := range stockCodes {
	//	fmt.Println(c)
	//}
	stockCodes = ds.GetStockCodes("sh")
	fmt.Println(len(stockCodes))
	//for _, c := range stockCodes {
	//	fmt.Println(c)
	//}
}

func TestTdxDataSource_GetStockNameHistory(t *testing.T) {
	ds := tdxdatasource.NewDataSource("data", true)
	codes := append(ds.GetStockCodes("sz"), ds.GetStockCodes("sh")...)

	for _, code := range codes[:1] {
		security, _ := entity.ParseSecurity(code)
		items := ds.GetStockNameHistory(security)
		fmt.Printf("%s, %+v\n", code, items)
	}
}

func TestTdxDataSource_GetStockName(t *testing.T) {
	ds := tdxdatasource.NewDataSource("data", true)
	codes := ds.GetStockCodes("sz")

	for _, code := range codes {
		security, _ := entity.ParseSecurity(code)
		name := ds.GetStockName(security)
		fmt.Printf("%s, %+v\n", code, name)
	}
}
