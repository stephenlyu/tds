package secondrecordgenerator

import (
	"testing"
	"github.com/stephenlyu/tds/entity"
	"io/ioutil"
	"github.com/stephenlyu/tds/util"
	"strings"
	"encoding/json"
	"fmt"
)

func loadTestData() []entity.TickItem {
	raw, err := ioutil.ReadFile("../20181009/EOSQFUT.OKEX")
	util.Assert(err == nil, "")

	var ret []entity.TickItem

	lines := strings.Split(string(raw), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var tick entity.TickItem
		err = json.Unmarshal([]byte(line), &tick)
		util.Assert(err == nil, "")

		ret = append(ret, tick)
	}
	return ret
}

func TestRecordGenerator_Feed(t *testing.T) {
	ticks := loadTestData()

	security := entity.ParseSecurityUnsafe("EOSQFUT.OKEX")

	rg := NewSecondRecordGenerator(security)

	for i := range ticks {
		r := rg.Feed(&ticks[i])
		if r == nil {
			continue
		}
		fmt.Printf("%+v\n", r)
	}
}
