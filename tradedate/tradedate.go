package tradedate

import (
	"fmt"
	"sort"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/stephenlyu/tds"
	"github.com/stephenlyu/tds/date"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/util"
)

const (
	TRADE_DATE_SEP_TIME  = "17:00:00"
	TRADE_DATE_DAY_START = "08:00:00"

	DAY_MILLIS      = 24 * 60 * 60 * 1000
	MINUTE_MILLIS   = 60 * 1000
	MINUTES_PER_DAY = DAY_MILLIS / MINUTE_MILLIS
)

type _TickerCacheItem struct {
	startTs, endTs uint64
	Tickers        []uint64
}

type _TradeDateCache map[string][]string
type _TickerCache map[string]*_TickerCacheItem

var tradeDateCache = make(_TradeDateCache) // key: exchange  	value: list of trade date
var tickerCache = make(_TickerCache)
var tradeDateCacheLock sync.RWMutex

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

func (this *_TickerCache) Set(commCode string, item *_TickerCacheItem) {
	tradeDateCacheLock.Lock()
	defer tradeDateCacheLock.Unlock()
	(*this)[commCode] = item
}

func (this *_TickerCache) Get(commCode string) *_TickerCacheItem {
	tradeDateCacheLock.RLock()
	defer tradeDateCacheLock.RUnlock()

	r, _ := (*this)[commCode]
	return r
}

func GetSecurityTradeDates(security *entity.Security) []string {
	ex := security.GetExchange()
	ret := tradeDateCache.Get(ex)
	if ret != nil {
		return ret
	}

	meta := TRADE_META.GetTradeDateMeta(ex)
	weekendTrading := TRADE_META.IsWeekendTrading(ex)

	startTs, err := date.DayString2Timestamp(meta.From)
	util.Assert(err == nil, "")

	endTs, err := date.DayString2Timestamp(meta.To)
	util.Assert(err == nil, "")

	nonDatesSet := mapset.NewSet()
	for _, d := range meta.NonTradingDates {
		nonDatesSet.Add(d)
	}

	for ts := startTs; ts <= endTs; ts += DAY_MILLIS {
		d := time.Unix(int64(ts)/1000, (int64(ts)%1000)*int64(time.Millisecond)).In(tds.Local)
		ds := d.Format(date.DAY_FORMAT)
		if nonDatesSet.Contains(ds) {
			continue
		}

		if !weekendTrading && (d.Weekday() == time.Saturday || d.Weekday() == time.Sunday) {
			continue
		}
		ret = append(ret, ds)
	}

	tradeDateCache.Set(ex, ret)
	return ret
}

func GetTradeDateRangeByDateString(security *entity.Security, dateString string) (startTs string, endTs string, thisTradeDate string, prevTradeDate string) {
	dates := GetSecurityTradeDates(security)

	for i := 1; i < len(dates); i++ {
		d := dates[i]
		pd := dates[i-1]
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
	endTs = date.Timestamp2SecondString(ts + 24*60*60*1000)

	return
}

func ToTradeTicker(security *entity.Security, timestamp uint64) uint64 {
	tickers := GetTradeTickers(security, timestamp)
	if timestamp < tickers[0] {
		return tickers[0]
	}
	lastTicker := tickers[len(tickers)-1]
	if timestamp >= lastTicker {
		return lastTicker
	}

	for i, ticker := range tickers {
		if i == 0 {
			continue
		}
		if timestamp < ticker {
			return tickers[i-1]
		}
	}

	util.UnreachableCode()
	return 0
}

func GetTradeTickers(security *entity.Security, timestamp uint64) []uint64 {
	commCode := fmt.Sprintf("%s.%s", security.Category, security.Exchange)
	item := tickerCache.Get(commCode)
	if item != nil && item.startTs <= timestamp && timestamp < item.endTs {
		return item.Tickers
	}

	dateString := date.Timestamp2SecondString(timestamp)
	startTs, endTs, tradeDate, preTradeDate := GetTradeDateRangeByDateString(security, dateString)

	isNonNight := TRADE_META.IsNonNightDate(security.Exchange, tradeDate)

	rawSpans := TRADE_META.GetDateTimeSpans(security, tradeDate)
	// 过滤夜盘的Spans
	var spans []TimeSpan
	if !isNonNight {
		spans = rawSpans
	} else {
		for _, span := range rawSpans {
			if span.Start >= TRADE_DATE_DAY_START && span.Start < TRADE_DATE_SEP_TIME {
				spans = append(spans, span)
			}
		}
	}

	// 根据Spans计算Tickers

	var tickers []uint64

	for _, span := range spans {
		start, end := span.Start, span.End
		if end == "24:00:00" {
			end = "23:59:59"
		}

		var from, to string
		if start >= TRADE_DATE_SEP_TIME {
			from = fmt.Sprintf("%s %s", preTradeDate, start)
			to = fmt.Sprintf("%s %s", preTradeDate, end)
		} else {
			from = fmt.Sprintf("%s %s", tradeDate, start)
			to = fmt.Sprintf("%s %s", tradeDate, end)
		}

		fromTs, _ := date.SecondString2Timestamp(from)
		toTs, _ := date.SecondString2Timestamp(to)
		for ts := fromTs; ts < toTs; ts += MINUTE_MILLIS {
			tickers = append(tickers, ts)
		}
	}

	sort.SliceStable(tickers, func(i, j int) bool {
		return tickers[i] < tickers[j]
	})

	start, _ := date.SecondString2Timestamp(startTs)
	end, _ := date.SecondString2Timestamp(endTs)
	tickerCache.Set(commCode, &_TickerCacheItem{
		startTs: start,
		endTs:   end,
		Tickers: tickers,
	})

	return tickers
}

// @param now - current time
// @param tradeTime - trade time span
// @param delayMinutes - munutes extended from end of timespan
func IsInTimeRange(now uint64, timeranges [][2]string, delayMinutes int64) bool {
	delayMillis := delayMinutes * 60 * 1000
	today := date.Timestamp2DayString(now)[:8]
	for _, span := range timeranges {
		start := fmt.Sprintf("%s %s", today, span[0])
		end := fmt.Sprintf("%s %s", today, span[1])
		startTs, _ := date.SecondString2Timestamp(start)
		endTs, _ := date.SecondString2Timestamp(end)

		if now >= startTs && now <= uint64(int64(endTs)+delayMillis) {
			return true
		}
	}
	return false
}
