package secondsmoother

import (
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/datahandler"
	"github.com/stephenlyu/tds/util"
	"github.com/stephenlyu/tds/tradedate"
	"github.com/stephenlyu/tds/date"
)

//
// 1分钟K线平滑程序，根据交易时刻表，使用前一根K线的收盘价生成K线，并补足缺失的K线
//

const SECOND_MILLIS = 1000

type _SecondSmoother struct {
	security *entity.Security
	prevRecord *entity.Record

	tickers []uint64
	currentIndex int
	startTs, endTs uint64
}

func NewSecondSmoother(security *entity.Security, initPrevRecord *entity.Record) datahandler.RecordHandler {
	util.Assert(security != nil, "")

	ret := &_SecondSmoother{
		security: security,
		prevRecord: initPrevRecord,
	}

	if initPrevRecord != nil {
		ret.init(ret.prevRecord)
		ret.calcCurrentIndex()
	}

	return ret
}

func (this *_SecondSmoother) init(r *entity.Record) {
	startTs, endTs, _, _ := tradedate.GetTradeDateRangeByDateString(this.security, r.GetDate())
	minutes := tradedate.GetTradeTickers(this.security, r.Date)

	tickers := make([]uint64, len(minutes) * 60)
	for i, m := range minutes {
		for j := 0; j < 60; j++ {
			tickers[i * 60 + j] = m + uint64(j) * 1000
		}
	}
	this.tickers = tickers

	this.startTs, _ = date.SecondString2Timestamp(startTs)
	this.endTs, _ = date.SecondString2Timestamp(endTs)
}

func (this *_SecondSmoother) calcCurrentIndex() {
	this.currentIndex = util.FindUInt64s(this.tickers, this.prevRecord.Date)
	util.Assert(this.currentIndex != -1, "")
}

func (this *_SecondSmoother) Feed(r *entity.Record) []*entity.Record {
	if this.prevRecord == nil {
		this.prevRecord = r
		this.init(r)
		this.calcCurrentIndex()
		return []*entity.Record{r}
	}

	switch {
	case r.Date >= this.startTs && r.Date < this.endTs:
	default:
		// 跨越交易日
		this.init(r)
		this.currentIndex = -1
	}

	index := util.FindUInt64s(this.tickers, r.Date)
	util.Assert(index != -1, "")

	var ret []*entity.Record
	// 每日第一根K，不需要平滑
	if index == 0 {
		goto done
	}

	// 休息时段不进行平滑
	if this.tickers[index] - this.tickers[index - 1] > SECOND_MILLIS {
		goto done
	}

	for i := this.currentIndex + 1; i < index; i++ {
		ret = append(ret, &entity.Record{
			Date: this.tickers[i],
			Open: this.prevRecord.Close,
			Close: this.prevRecord.Close,
			High: this.prevRecord.Close,
			Low: this.prevRecord.Close,
			Volume: 0,
			Amount: 0,
			BuyVolume: 0,
			SellVolume: 0,
		})
	}

done:
	ret = append(ret, r)
	this.prevRecord = r
	this.currentIndex = index
	return ret
}
