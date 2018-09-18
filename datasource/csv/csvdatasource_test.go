package csvdatasource_test

import (
	"testing"
	"github.com/stephenlyu/tds/datasource/tdx"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tds/entity"
	"fmt"
	"github.com/stephenlyu/tds/datasource/csv"
	"github.com/stephenlyu/tds/date"
)

func Test_CSVDataSource_SaveData(t *testing.T) {
	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	tdxDs := tdxdatasource.NewDataSource("../tdx/data", true)

	err, data := tdxDs.GetData(security, period1)
	util.Assert(err == nil, "")

	csvDs := csvdatasource.NewCSVDataSource("csv")
	err = csvDs.SaveData(security, period1, data[:1000])
	util.Assert(err == nil, fmt.Sprintf("%+v", err))
}

func Test_CSVDataSource_GetData(t *testing.T) {
	csvDs := csvdatasource.NewCSVDataSource("csv")
	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	start := util.Tick()
	err, data := csvDs.GetData(security, period1)
	fmt.Printf("time cost: %dms\n", util.Tick() - start)
	util.Assert(err == nil, "")
	fmt.Println(len(data))
	fmt.Printf("%+v\n", &data[0])
	fmt.Printf("%+v\n", &data[len(data) - 1])
}

func Test_CSVDataSource_GetRangeData(t *testing.T) {
	csvDs := csvdatasource.NewCSVDataSource("csv")
	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	startDate, _ := date.SecondString2Timestamp("20150217 14:31:00")
	err, data := csvDs.GetRangeData(security, period1, startDate, startDate)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}
}
