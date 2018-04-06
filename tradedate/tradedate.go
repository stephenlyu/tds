package tradedate

import (
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tds/date"
)

func GetTradeDateRange(security *entity.Security, dateString string) (startTs string, endTs string) {
	util.Assert(len(dateString) >= 8, "")
	day := dateString[:8]

	startTs = day + " 00:00:00"
	ts, _ := date.SecondString2Timestamp(startTs)
	endTs = date.Timestamp2SecondString(ts + 24 * 60 * 60 * 1000)

	return
}
