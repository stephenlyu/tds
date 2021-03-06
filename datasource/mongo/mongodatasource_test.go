package mongodatasource_test

import (
	"testing"
	"github.com/stephenlyu/tds/datasource/tdx"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/datasource/mongo"
	"fmt"
)

func Test_MongoDataSource_SaveData(t *testing.T) {
	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	tdxDs := tdxdatasource.NewDataSource("../tdx/data", true)

	err, data := tdxDs.GetData(security, period1)
	util.Assert(err == nil, "")

	mongoDs := mongodatasource.NewMongoDataSource("localhost", "data")
	err = mongoDs.AppendData(security, period1, data[:1000])
	util.Assert(err == nil, "")
}

func Test_MongoDataSource_GetData(t *testing.T) {
	mongoDs := mongodatasource.NewMongoDataSource("localhost", "data")
	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	start := util.Tick()
	err, data := mongoDs.GetData(security, period1)
	fmt.Printf("time cost: %dms\n", util.Tick() - start)
	util.Assert(err == nil, "")
	fmt.Println(len(data))
	fmt.Printf("%+v\n", &data[0])
	fmt.Printf("%+v\n", &data[len(data) - 1])


	err, r := mongoDs.GetLastRecord(security, period1)
	util.Assert(err == nil, fmt.Sprintf("%+v", err))
	fmt.Printf("%+v\n", r)

	err, data = mongoDs.GetDataEx(security, period1, 1423704660000, 100)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 100, "")

	err = mongoDs.RemoveData(security, period1, 0, 0)
	util.Assert(err == nil, "")
	err, data = mongoDs.GetData(security, period1)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 0, "")
}