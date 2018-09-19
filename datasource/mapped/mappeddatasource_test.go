package mappeddatasource_test

import (
	"testing"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/util"
	. "github.com/stephenlyu/tds/entity"
	"fmt"
	"github.com/stephenlyu/tds/datasource/csv"
	"github.com/stephenlyu/tds/date"
	. "github.com/stephenlyu/tds/datasource"
	"github.com/stephenlyu/tds/datasource/mapped"
)

type _mapper struct {
}

func (this *_mapper) MapDateRanges(security *Security) []DateRange {
	sep, _ := date.SecondString2Timestamp("20180918 15:00:00")
	return []DateRange {
		{StartDate: 0, EndDate: sep, Security: ParseSecurityUnsafe("000001.SZ")},
		{StartDate: sep, EndDate: 0, Security: ParseSecurityUnsafe("000002.SZ")},
	}
}

type _nomapper struct {
}

func (this *_nomapper) MapDateRanges(security *Security) []DateRange {
	return []DateRange {
		{Security: security},
	}
}

func Test_MappedDataSource_GetData(t *testing.T) {
	csvDs := csvdatasource.NewCSVDataSource("csv")

	security, err := ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	mappedDs := mappeddatasource.NewMapperDataSource()
	mappedDs.SetMapper(&_mapper{})
	mappedDs.SetTargetDataSource(csvDs)

	err, data := mappedDs.GetData(security, period1)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}
}

func Test_MappedDataSource_GetRangeData(t *testing.T) {
	csvDs := csvdatasource.NewCSVDataSource("csv")

	security, err := ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	mappedDs := mappeddatasource.NewMapperDataSource()
	mappedDs.SetMapper(&_mapper{})
	mappedDs.SetTargetDataSource(csvDs)

	startDate, _ := date.SecondString2Timestamp("20180918 14:30:00")
	endDate, _ := date.SecondString2Timestamp("20180919 10:00:00")

	err, data := mappedDs.GetRangeData(security, period1, startDate, endDate)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}

	fmt.Println("====================================")

	endDate1, _ := date.SecondString2Timestamp("20180918 15:00:00")

	err, data = mappedDs.GetRangeData(security, period1, startDate, endDate1)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}

	fmt.Println("====================================")

	startDate1, _ := date.SecondString2Timestamp("20180918 15:00:00")

	err, data = mappedDs.GetRangeData(security, period1, startDate1, 0)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}

	fmt.Println("====================================")

	err, data = mappedDs.GetRangeData(security, period1, 0, 0)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}
}

func Test_MappedDataSource_GetData_Nomapper(t *testing.T) {
	csvDs := csvdatasource.NewCSVDataSource("csv")

	security, err := ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	mappedDs := mappeddatasource.NewMapperDataSource()
	mappedDs.SetMapper(&_nomapper{})
	mappedDs.SetTargetDataSource(csvDs)

	err, data := mappedDs.GetData(security, period1)
	util.Assert(err == nil, "")
	for i := range data {
		fmt.Printf("%+v\n", &data[i])
	}
}
