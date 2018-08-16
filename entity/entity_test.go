package entity_test

import (
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/util"
	"testing"
)

func TestSecurity(t *testing.T) {
	code := "600000.SH"
	security, err := entity.ParseSecurity(code)
	util.Assert(err == nil, "err == nil")
	util.Assert(security != nil, "")
	util.Assert(security.String() == code, "")
}
