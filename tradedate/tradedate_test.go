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
	tickers := GetTradeTickers(security, 1534671000000)
	for i, ticker := range tickers {
		fmt.Println(i, date.Timestamp2SecondString(ticker))
	}
}

func TestGetTradeDateRangeByDateString(t *testing.T) {
	security := entity.ParseSecurityUnsafe("EOSQFUT.OKEX")
	today := date.GetTodayString()

	startTs, endTs, _, _ := GetTradeDateRangeByDateString(security, today)
	fmt.Println(startTs, endTs)
}

func TestGetNonNightTradeTickers(t *testing.T) {
	security := entity.ParseSecurityUnsafe("J1901.DCE")
	ts, _ := date.DayString2Timestamp("20181008")
	tickers := GetTradeTickers(security, ts)
	for i, ticker := range tickers {
		fmt.Println(i, date.Timestamp2SecondString(ticker))
	}
}

func TestToTradeTicker(t *testing.T) {
	security := entity.ParseSecurityUnsafe("J1901.DCE")
	ts, _ := date.SecondString2Timestamp("20181008 15:03:00")
	ticker := ToTradeTicker(security, ts)
	fmt.Println(date.Timestamp2SecondString(ticker))
}
