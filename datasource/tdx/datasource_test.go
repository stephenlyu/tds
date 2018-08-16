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
	"os"
	"github.com/stephenlyu/tds/datasource"
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

	err, data = ds.GetDataEx(security, period2, startDate, 100)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 100, "")

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
	codes := append(ds.GetStockCodes("sz"), ds.GetStockCodes("sh")...)

	for _, code := range codes {
		security, _ := entity.ParseSecurity(code)
		name := ds.GetStockName(security)
		fmt.Printf("%s, %+v\n", code, name)
	}
}

func TestTdxDataSource_GetStockNames(t *testing.T) {
	ds := tdxdatasource.NewDataSource("data", true)

	names := ds.GetStockNames()

	for k, v := range names {
		fmt.Printf("%s, %+v\n", k, v)
	}
}

func TestTdxDataSource_GetLastRecord(t *testing.T) {
	ds := tdxdatasource.NewDataSource("data", true)
	security, _ := entity.ParseSecurity("000001.SZ")
	_, period := period.PeriodFromString("D1")
	err, r := ds.GetLastRecord(security, period)
	util.Assert(err == nil, "")
	fmt.Printf("%+v\n", r)
}

var eqRecords = func(rs1, rs2 []entity.Record) bool {
	if len(rs1) != len(rs2) {
		return false
	}

	for i := range rs1 {
		if !rs1[i].Eq(&rs2[i]) {
			return false
		}
	}
	return true
}

func TestTdxDataSource_AppendData(t *testing.T) {
	ds := tdxdatasource.NewDataSource("data", true)
	os.RemoveAll("temp")

	ds1 := tdxdatasource.NewDataSource("temp", true)

	security, _ := entity.ParseSecurity("000001.SZ")
	_, period := period.PeriodFromString("D1")

	err, records := ds.GetData(security, period)
	util.Assert(err == nil, "")
	fmt.Printf("record count: %d\n", len(records))

	err = ds1.AppendData(security, period, records)
	util.Assert(err == nil, fmt.Sprintf("%v", err))

	_, records1 := ds1.GetData(security, period)
	util.Assert(eqRecords(records, records1), "")

}

func TestCustomPeriodSave(t *testing.T) {
	ds := tdxdatasource.NewDataSource("data", true)
	security, _ := entity.ParseSecurity("000001.SZ")
	err, records := ds.GetData(security, period.PERIOD_M5)
	util.Assert(err == nil, "")
	fmt.Printf("%s\n", records[len(records) - 1].String())

	converter := datasource.NewPeriodConverter(period.PERIOD_M5, period.PERIOD_M15)
	destData := converter.Convert(records)
	ds.SaveData(security, period.PERIOD_M15, destData)

	err, data := ds.GetData(security, period.PERIOD_M15)
	util.Assert(err == nil, "")
	util.Assert(len(data) == len(destData), "")
	for i := range destData {
		r1 := &destData[i]
		r2 := &data[i]

		util.Assert(r1.GetDate() == r2.GetDate(), "")
		util.Assert(util.Round(float64(r1.GetOpen()), 2) == util.Round(float64(r2.GetOpen()), 2), fmt.Sprintf("%s %d %d", r1.GetDate(), r1.Open, r2.Open))
		util.Assert(util.Round(float64(r1.GetClose()), 2) == util.Round(float64(r2.GetClose()), 2), fmt.Sprintf("%s %d %d", r1.GetDate(), r1.Close, r2.Close))
		util.Assert(util.Round(float64(r1.GetLow()), 2) == util.Round(float64(r2.GetLow()), 2), fmt.Sprintf("%s %d %d", r1.GetDate(), r1.Low, r2.Low))
		util.Assert(util.Round(float64(r1.GetHigh()), 2) == util.Round(float64(r2.GetHigh()), 2), fmt.Sprintf("%s %d %d", r1.GetDate(), r1.High, r2.High))
	}
}

func TestTdxDataSource_AppendRawData(t *testing.T) {
	ds := tdxdatasource.NewDataSource("data", true)
	os.RemoveAll("temp")

	ds1 := tdxdatasource.NewDataSource("temp", true)

	security, _ := entity.ParseSecurity("000001.SZ")
	_, period := period.PeriodFromString("D1")

	err, records := ds.GetData(security, period)
	util.Assert(err == nil, "")
	fmt.Printf("record count: %d\n", len(records))

	marshaller := tdxdatasource.NewMarshaller(period)
	raw := make([]byte, len(records) * tdxdatasource.TDX_RECORD_SIZE)
	for i := range records {
		start := i * tdxdatasource.TDX_RECORD_SIZE
		end := start + tdxdatasource.TDX_RECORD_SIZE
		bytes, _ := marshaller.ToBytes(&records[i])
		copy(raw[start:end], bytes)
	}

	// Append All
	err = ds1.AppendRawData(security, period, raw)
	util.Assert(err == nil, fmt.Sprintf("%v", err))

	_, records1 := ds1.GetData(security, period)
	util.Assert(eqRecords(records, records1), "")

	// Overlap case
	err = ds1.AppendRawData(security, period, raw[:1000*tdxdatasource.TDX_RECORD_SIZE])
	util.Assert(err == nil, fmt.Sprintf("%v", err))

	_, records1 = ds1.GetData(security, period)
	util.Assert(eqRecords(records[:1000], records1), "")

	err = ds1.AppendRawData(security, period, raw[900*tdxdatasource.TDX_RECORD_SIZE:])
	util.Assert(err == nil, fmt.Sprintf("%v", err))

	_, records1 = ds1.GetData(security, period)
	util.Assert(eqRecords(records, records1), "")

	// Normal Case
	err = ds1.AppendRawData(security, period, raw[:1000*tdxdatasource.TDX_RECORD_SIZE])
	util.Assert(err == nil, fmt.Sprintf("%v", err))

	_, records1 = ds1.GetData(security, period)
	util.Assert(eqRecords(records[:1000], records1), "")

	err = ds1.AppendRawData(security, period, raw[1000*tdxdatasource.TDX_RECORD_SIZE:])
	util.Assert(err == nil, fmt.Sprintf("%v", err))

	_, records1 = ds1.GetData(security, period)
	util.Assert(eqRecords(records, records1), "")
}