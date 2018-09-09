package m1smoother

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

const MINUTE_MILLIS = 60 * 1000

type _M1Smoother struct {
	security *entity.Security
	prevRecord *entity.Record

	tickers []uint64
	currentIndex int
	startTs, endTs uint64
}

func NewM1Smoother(security *entity.Security, initPrevRecord *entity.Record) datahandler.RecordHandler {
	util.Assert(security != nil, "")
	util.Assert(initPrevRecord != nil, "")

	ret := &_M1Smoother{
		security: security,
		prevRecord: initPrevRecord,
	}

	ret.init(ret.prevRecord)
	ret.calcCurrentIndex()

	return ret
}

func (this *_M1Smoother) init(r *entity.Record) {
	startTs, endTs, _, _ := tradedate.GetTradeDateRangeByDateString(this.security, r.GetDate())
	this.tickers = tradedate.GetTradeTickers(this.security, r.Date)

	this.startTs, _ = date.SecondString2Timestamp(startTs)
	this.endTs, _ = date.SecondString2Timestamp(endTs)
}

func (this *_M1Smoother) calcCurrentIndex() {
	this.currentIndex = util.FindUInt64s(this.tickers, this.prevRecord.Date)
	util.Assert(this.currentIndex != -1, "")
}

func (this *_M1Smoother) Feed(r *entity.Record) []*entity.Record {
	util.Assert(this.prevRecord != nil, "")

	var ret []*entity.Record
	switch {
	case r.Date >= this.startTs && r.Date < this.endTs:
	default:
		// 前一日尾部缺数据，需要平滑出来
		for i := this.currentIndex + 1; i < len(this.tickers); i++ {
			ret = append(ret, &entity.Record{
				Date: this.tickers[i],
				Open: this.prevRecord.Close,
				Close: this.prevRecord.Close,
				High: this.prevRecord.Close,
				Low: this.prevRecord.Close,
				Volume: 0,
				Amount: 0,
			})
		}

		// 跨越交易日
		this.init(r)
		this.currentIndex = -1
	}

	index := util.FindUInt64s(this.tickers, r.Date)
	util.Assert(index != -1, "")

	for i := this.currentIndex + 1; i < index; i++ {
		ret = append(ret, &entity.Record{
			Date: this.tickers[i],
			Open: this.prevRecord.Close,
			Close: this.prevRecord.Close,
			High: this.prevRecord.Close,
			Low: this.prevRecord.Close,
			Volume: 0,
			Amount: 0,
		})
	}

	ret = append(ret, r)

	this.prevRecord = r
	this.currentIndex = index
	return ret
}
