package tradedate

import (
	"io/ioutil"
	"log"
	"encoding/json"
	"sync"
	"github.com/stephenlyu/tds/entity"
	"fmt"
	"github.com/stephenlyu/tds/util"
	"github.com/sirupsen/logrus"
)

type TimeSpan struct {
	Start, End string
}

type TimeSpans struct {
	From string
	To string
	Spans []TimeSpan
}

type TradeTimeSpanDesc struct {
	Spans []TimeSpans
}

type TradeDateMeta struct {
	From string
	To string
	NonTradingDates []string
}

type TradeMeta struct {
	WeekendTradingExchanges map[string]bool
	TradingTimeMeta map[string][]TimeSpans
	NonNightDates map[string][]string
	TradeDateMeta map[string]*TradeDateMeta
}

func (this *TradeMeta) IsWeekendTrading(exchange string) bool {
	ret, ok := this.WeekendTradingExchanges[exchange]
	if !ok {
		return false
	}
	return ret
}

func (this *TradeMeta) GetTimeSpans(security *entity.Security) []TimeSpans {
	commCode := fmt.Sprintf("%s.%s", security.Category, security.Exchange)
	ret, ok := this.TradingTimeMeta[commCode]
	if ok {
		return ret
	}

	ret, ok = this.TradingTimeMeta[security.Exchange]
	return ret
}

func (this *TradeMeta) GetDateTimeSpans(security *entity.Security, date string) []TimeSpan {
	spans := this.GetTimeSpans(security)
	for _, o := range spans {
		if o.From <= date && date < o.To {
			return o.Spans
		}
	}

	util.UnreachableCode()
	return nil
}

func (this *TradeMeta) IsNonNightDate(exchange string, d string) bool {
	ret, ok := this.NonNightDates[exchange]
	if !ok {
		return false
	}

	for _, s := range ret {
		if s == d {
			return true
		}
	}

	return false
}

func (this *TradeMeta) GetTradeDateMeta(exchange string) *TradeDateMeta {
	ret, _ := this.TradeDateMeta[exchange]
	return ret
}

var TRADE_META *TradeMeta = &TradeMeta{
	WeekendTradingExchanges: map[string]bool {
		"OKEX": true,
		"BITMEX": true,
	},
	TradingTimeMeta: map[string][]TimeSpans {
		"OKEX": []TimeSpans{
			{
				From: "20130706",
				To: "20380101",
				Spans: []TimeSpan{
					{
						Start: "00:00:00",
						End: "17:00:00",
					},
					{
						Start: "17:00:00",
						End: "24:00:00",
					},
				},
			},
		},
		"BITMEX": []TimeSpans{
			{
				From: "20130706",
				To: "20380101",
				Spans: []TimeSpan{
					{
						Start: "00:00:00",
						End: "17:00:00",
					},
					{
						Start: "17:00:00",
						End: "24:00:00",
					},
				},
			},
		},
	},
	NonNightDates: map[string][]string{
	},
	TradeDateMeta: map[string]*TradeDateMeta {
		"OKEX": {
			From: "20150101",
			To: "20381230",
			NonTradingDates: []string{
			},
		},
		"BITMEX": {
			From: "20150101",
			To: "20381230",
			NonTradingDates: []string{
			},
		},
	},
}

func LoadMeta() {
	once := sync.Once{}
	once.Do(
		func() {
			bytes, err := ioutil.ReadFile("trade-meta.json")
			if err != nil {
				logrus.Errorf("Load trade-meta.json fail. error: %+v", err)
				return
			}

			var meta *TradeMeta
			err = json.Unmarshal(bytes, &meta)
			if err != nil {
				log.Fatal(err)
			}
			TRADE_META = meta
		})
}

