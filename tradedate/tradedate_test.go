package tradedate

import (
	"testing"
	"fmt"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tds/date"
)

func TestGetTradeDateRange(t *testing.T) {
	startTs, endTs := GetTradeDateRange(nil, "20180405")
	fmt.Println(startTs, endTs)
}

func TestToTradeMinute(t *testing.T) {
	security := entity.ParseSecurityUnsafe("EOSQFUT.OKEX")
	fmt.Println(date.Timestamp2SecondString(ToTradeTicker(security, util.Tick())))

	ts, _ := date.SecondString2Timestamp("20180819 11:23:00")
	fmt.Println(date.Timestamp2SecondString(ToTradeTicker(security, ts)))

	ts, _ = date.SecondString2Timestamp("20180819 11:23:01")
	fmt.Println(date.Timestamp2SecondString(ToTradeTicker(security, ts)))

	ts, _ = date.SecondString2Timestamp("20180819 11:23:59")
	fmt.Println(date.Timestamp2SecondString(ToTradeTicker(security, ts)))
}

func TestGetTradeTickers(t *testing.T) {
	security := entity.ParseSecurityUnsafe("EOSQFUT.OKEX")
	tickers := GetTradeTickers(security, util.Tick())
	for i, ticker := range tickers {
		fmt.Println(i, date.Timestamp2SecondString(ticker))
	}
}