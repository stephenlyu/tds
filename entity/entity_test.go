package entity_test

import (
	"github.com/stephenlyu/tds/entity"
	"time"
	"github.com/stephenlyu/tds/util"
	"testing"
)

func TestEntity(t *testing.T) {
	r := &entity.Record{
		Date: uint64(time.Now().UnixNano() / int64(time.Millisecond)),
		Open: 8500,
		Close: 9050,
		High: 9190,
		Low: 8460,
		Volume: 10000,
		Amount: 1000000,
	}

	r1 := &entity.Record{}
	err := entity.RecordFromBytes(r.Bytes(), r1)
	util.Assert(err == nil, "")
	util.Assert(r.Eq(r1), "r.Eq(r1)")
}

func TestSecurity(t *testing.T) {
	code := "600000.SH"
	security, err := entity.ParseSecurity(code)
	util.Assert(err == nil, "err == nil")
	util.Assert(security != nil, "")
	util.Assert(security.String() == code, "")
}
