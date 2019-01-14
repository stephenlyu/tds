package entity_test

import (
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/util"
	"testing"
	"fmt"
	"encoding/json"
)

func TestSecurity(t *testing.T) {
	codes := []string{
		"600000.SH", "EOSQFUT.OKEX", "BTCFUT.BITMEX", "EOSFUT.PLO",
		"BTCUSDTSPOT.HUOBI", "RB1901.SHFE",
	}

	for _, code := range codes {
		security, err := entity.ParseSecurity(code)
		util.Assert(err == nil, fmt.Sprintf("%+v", err))
		util.Assert(security != nil, "")
		util.Assert(security.String() == code, code)
		fmt.Printf("cat: %s code: %s exchange: %s\n", security.GetCategory(), security.GetCode(), security.GetExchange())
	}
}

func TestRecord_ToProtoBytes(t *testing.T) {
	record := entity.Record{Date:1483407240000, Open:9.11, Close:9.099, High:9.11, Low:9.09, Volume:2.0228e+06, Amount:1.8414828e+07}

	start := util.NanoTick()
	bytes, err := record.ToProtoBytes()
	fmt.Printf("proto marshal cost: %dns\n", util.NanoTick() - start)
	util.Assert(err == nil, "")
	fmt.Printf("proto byte len: %d\n", len(bytes))

	zbytes := util.ZlibCompress(bytes)
	fmt.Printf("compressed proto byte len: %d\n", len(zbytes))

	start = util.NanoTick()
	r, err := entity.RecordFromProtoBytes(bytes)
	fmt.Printf("proto unmarshal cost: %dns\n", util.NanoTick() - start)
	util.Assert(err == nil, "")
	util.Assert(r.Eq(&record), "")

	start = util.NanoTick()
	bytes, _ = json.Marshal(record)
	fmt.Printf("json marshal cost: %dns\n", util.NanoTick() - start)
	fmt.Printf("json byte len: %d\n", len(bytes))

	var nr entity.Record
	json.Unmarshal(bytes, &nr)
	fmt.Printf("json unmarshal cost: %dns\n", util.NanoTick() - start)
}
