package tradedate

import (
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tds/date"
	"sync"
)

// TODO:

const (
	TRADE_DATE_SEP_TIME = "17:00:00"
	TRADE_DATE_DAY_START = "08:00:00"

	DAY_MILLIS = 24 * 60 * 60 * 1000
	MINUTE_MILLIS = 60 * 1000
	MINUTES_PER_DAY = DAY_MILLIS / MINUTE_MILLIS
)

type TimeSpan struct {
	Start, End string
}

type TradeTimeSpanDesc struct {
	Spans []TimeSpan
}

type TradeDateDesc struct {
	StartDate string
	WeekendTrading bool						// true 表示周末交易
	NonTradeDates []string
	DefaultTimeSpanDesc *TradeTimeSpanDesc
}

var ALL_DAY_TRADE_TIME_SPAN_DESC = &TradeTimeSpanDesc{
}

var globalTradeDateDesc = map[string]*TradeDateDesc {
	"OKEX": &TradeDateDesc{
		StartDate: "20140101",
		WeekendTrading: true,
		NonTradeDates: []string{},
		DefaultTimeSpanDesc: ALL_DAY_TRADE_TIME_SPAN_DESC,
	},
}

var globalTradeTimeSpanDesc = map[string]*TradeTimeSpanDesc {
	// TODO:
}

type _TradeDateCache map[string][]string

var tradeDateCache = make(_TradeDateCache) // key: exchange  	value: list of trade date
var tradeDateCacheLock sync.RWMutex

func (this *TradeTimeSpanDesc) isAllDay() bool {
	return this == ALL_DAY_TRADE_TIME_SPAN_DESC
}

func (this *_TradeDateCache) Set(exchange string, dates []string) {
	tradeDateCacheLock.Lock()
	defer tradeDateCacheLock.Unlock()
	tradeDateCache[exchange] = dates
}

func (this *_TradeDateCache) Get(exchange string) []string {
	tradeDateCacheLock.RLock()
	defer tradeDateCacheLock.RUnlock()

	r, _ := tradeDateCache[exchange]
	return r
}

func GetSecurityTradeDates(security *entity.Security) []string {
	ex := security.GetExchange()
	ret := tradeDateCache.Get(ex)
	if ret != nil {
		return ret
	}

	dd, ok := globalTradeDateDesc[ex]
	util.Assert(ok, "")

	startTs, err := date.DayString2Timestamp(dd.StartDate)
	util.Assert(err == nil, "")

	ret = nil
	for v := startTs; v < util.Tick(); v += DAY_MILLIS {
		d := date.Timestamp2DayString(v)
		if util.InStrings(d, dd.NonTradeDates) {
			continue
		}
		ret = append(ret, d)
	}

	tradeDateCache.Set(ex, ret)
	return ret
}

func GetTradeDateRangeByDateString(security *entity.Security, dateString string) (startTs string, endTs string, thisTradeDate string, prevTradeDate string) {
	dates := GetSecurityTradeDates(security)

	for i := 1; i < len(dates); i++ {
		d := dates[i]
		pd := dates[i - 1]
		st := pd + " " + TRADE_DATE_SEP_TIME
		et := d + " " + TRADE_DATE_SEP_TIME
		if dateString >= st && dateString < et {
			startTs = st
			endTs = et
			thisTradeDate = d
			prevTradeDate = pd
			return
		}
	}
	return
}

// FIXME:
func GetTradeDateRange(security *entity.Security, dateString string) (startTs string, endTs string) {
	util.Assert(len(dateString) >= 8, "")
	day := dateString[:8]

	startTs = day + " 00:00:00"
	ts, _ := date.SecondString2Timestamp(startTs)
	endTs = date.Timestamp2SecondString(ts + 24 * 60 * 60 * 1000)

	return
}

func getTradeTimeSpanDesc(security *entity.Security) *TradeTimeSpanDesc {
	if s, ok := globalTradeTimeSpanDesc[security.GetExchange()]; ok {
		return s
	}

	if dd, ok := globalTradeDateDesc[security.GetExchange()]; ok {
		return dd.DefaultTimeSpanDesc
	}

	util.UnreachableCode()
	return nil
}

func ToTradeTicker(security *entity.Security, timestamp uint64) uint64 {
	tsd := getTradeTimeSpanDesc(security)

	if tsd.isAllDay() {
		timestamp = timestamp / MINUTE_MILLIS * MINUTE_MILLIS
		return timestamp
	}

	util.UnreachableCode()
	return 0
}

func GetTradeTickers(security *entity.Security, timestamp uint64) []uint64 {
	tsd := getTradeTimeSpanDesc(security)
	dateString := date.Timestamp2SecondString(timestamp)
	startTs, endTs, _, _ := GetTradeDateRangeByDateString(security, dateString)
	if tsd.isAllDay() {
		st, _ := date.SecondString2Timestamp(startTs)
		et, _ := date.SecondString2Timestamp(endTs)

		ret := make([]uint64, MINUTES_PER_DAY)

		for i, ts := 0, st; ts < et; i, ts = i + 1, ts + MINUTE_MILLIS {
			ret[i] = ts
		}
		return ret
	}

	util.UnreachableCode()
	return nil
}
