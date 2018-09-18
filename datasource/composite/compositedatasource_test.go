package compositedatasource_test

import (
	"testing"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tds/entity"
	"fmt"
	"github.com/stephenlyu/tds/datasource/csv"
	"github.com/stephenlyu/tds/datasource/composite"
	"github.com/stephenlyu/tds/date"
)

func Test_CSVDataSource_GetData(t *testing.T) {
	csvDs := csvdatasource.NewCSVDataSource("csv")
	csvDs1 := csvdatasource.NewCSVDataSource("csv1")

	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	compositeDs := compositedatasource.NewCompositeDataSource()
	compositeDs.AddSubDatasource(csvDs)
	compositeDs.AddSubDatasource(csvDs1)

	err, data := compositeDs.GetData(security, period1)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}
}

func Test_CSVDataSource_GetDataEx(t *testing.T) {
	csvDs := csvdatasource.NewCSVDataSource("csv")
	csvDs1 := csvdatasource.NewCSVDataSource("csv1")

	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	compositeDs := compositedatasource.NewCompositeDataSource()
	compositeDs.AddSubDatasource(csvDs)
	compositeDs.AddSubDatasource(csvDs1)

	startDate, _ := date.SecondString2Timestamp("20150217 14:31:00")

	err, data := compositeDs.GetDataEx(security, period1, startDate, 200)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}
}

func Test_CSVDataSource_GetRangeData(t *testing.T) {
	csvDs := csvdatasource.NewCSVDataSource("csv")
	csvDs1 := csvdatasource.NewCSVDataSource("csv1")

	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	compositeDs := compositedatasource.NewCompositeDataSource()
	compositeDs.AddSubDatasource(csvDs)
	compositeDs.AddSubDatasource(csvDs1)

	startDate, _ := date.SecondString2Timestamp("20150217 14:31:00")
	endDate, _ := date.SecondString2Timestamp("20150217 14:57:00")

	err, data := compositeDs.GetRangeData(security, period1, startDate, endDate)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}

	fmt.Println("==========================================")

	endDate, _ = date.SecondString2Timestamp("20150225 09:40:00")

	err, data = compositeDs.GetRangeData(security, period1, startDate, endDate)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}
}

func Test_CSVDataSource_GetDataFromLast(t *testing.T) {
	csvDs := csvdatasource.NewCSVDataSource("csv")
	csvDs1 := csvdatasource.NewCSVDataSource("csv1")

	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	compositeDs := compositedatasource.NewCompositeDataSource()
	compositeDs.AddSubDatasource(csvDs)
	compositeDs.AddSubDatasource(csvDs1)

	err, data := compositeDs.GetDataFromLast(security, period1, 0, 2)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}

	fmt.Println("==========================================")

	endDate, _ := date.SecondString2Timestamp("20150225 09:40:00")
	err, data = compositeDs.GetDataFromLast(security, period1, endDate, 2)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}

	fmt.Println("==========================================")

	endDate, _ = date.SecondString2Timestamp("20150225 09:37:00")
	err, data = compositeDs.GetDataFromLast(security, period1, endDate, 2)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}

	fmt.Println("==========================================")

	endDate, _ = date.SecondString2Timestamp("20150217 14:59:00")
	err, data = compositeDs.GetDataFromLast(security, period1, endDate, 2)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}
}

func Test_CSVDataSource_GetLastRecord(t *testing.T) {
	csvDs := csvdatasource.NewCSVDataSource("csv")
	csvDs1 := csvdatasource.NewCSVDataSource("csv1")

	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	compositeDs := compositedatasource.NewCompositeDataSource()
	compositeDs.AddSubDatasource(csvDs)
	compositeDs.AddSubDatasource(csvDs1)

	err, lastR := compositeDs.GetLastRecord(security, period1)
	util.Assert(err == nil, "")
	fmt.Printf("%+v\n", lastR)
}