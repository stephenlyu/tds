package redisdatasource_test

import (
	"testing"
	"github.com/stephenlyu/tds/datasource/tdx"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tds/entity"
	"fmt"
	"github.com/stephenlyu/tds/datasource/redis"
	"math"
)

func Test_RedisDataSource_SaveData(t *testing.T) {
	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	tdxDs := tdxdatasource.NewDataSource("../tdx/data", true)

	err, data := tdxDs.GetData(security, period1)
	util.Assert(err == nil, "")

	redisDs := redisdatasource.NewRedisDataSource("", "")
	err = redisDs.AppendData(security, period1, data[:3000])
	util.Assert(err == nil, "")
}

func Test_RedisDataSource_GetData(t *testing.T) {
	redisDs := redisdatasource.NewRedisDataSource("", "")
	security, err := entity.ParseSecurity("000001.SZ")
	util.Assert(err == nil, "")
	util.Assert(security != nil, "")

	err, period1 := period.PeriodFromString("M1")
	util.Assert(err == nil, "")

	start := util.Tick()
	err, data := redisDs.GetData(security, period1)
	fmt.Printf("time cost: %dms\n", util.Tick() - start)
	util.Assert(err == nil, "")
	fmt.Println(len(data))
	fmt.Printf("%+v\n", &data[0])
	fmt.Printf("%+v\n", &data[len(data) - 1])

	err, data = redisDs.GetDataFromLast(security, period1, 0, 100)
	util.Assert(err == nil, "")
	fmt.Println(len(data))
	fmt.Printf("%+v\n", &data[0])
	fmt.Printf("%+v\n", &data[len(data) - 1])

	err, r := redisDs.GetLastRecord(security, period1)
	util.Assert(err == nil, fmt.Sprintf("%+v", err))
	fmt.Printf("%+v\n", r)

	err, data = redisDs.GetDataEx(security, period1, 1423704660000, 100)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 100, "")

	err = redisDs.RemoveData(security, period1, 0, math.MaxUint64)
	util.Assert(err == nil, "")
	err, data = redisDs.GetData(security, period1)
	util.Assert(err == nil, "")
	util.Assert(len(data) == 0, "")
}
