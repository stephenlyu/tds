package secondrecordgenerator

import (
	"github.com/stephenlyu/tds/entity"
	"math"
	"github.com/stephenlyu/tds/util"
)

//
// 从TickItem合成1分钟K线程序
//

type SecondRecordGenerator struct {
	Security *entity.Security

	Current *entity.Record		// 当前的K线
}

func NewSecondRecordGenerator(security *entity.Security) *SecondRecordGenerator {
	return &SecondRecordGenerator{
		Security: security,
	}
}

func (this *SecondRecordGenerator) Feed(tick *entity.TickItem) *entity.Record {
	util.Assert(tick.Code == this.Security.String(), "")
	if tick.Price == 0 || tick.Volume == 0 {
		return nil
	}

	ticker := tick.Timestamp / 1000 * 1000

	if this.Current == nil || ticker != this.Current.Date {
		this.Current = &entity.Record{
			Date: ticker,
			Open: tick.Price,
			Close: tick.Price,
			High: tick.High,
			Low: tick.Low,
			Volume: tick.Volume,
			Amount: tick.Amount,
		}

	} else {
		this.Current.Close = tick.Price
		this.Current.High = math.Max(this.Current.High, tick.High)
		this.Current.Low = math.Min(this.Current.Low, tick.Low)
		this.Current.Volume += tick.Volume
		this.Current.Amount += tick.Amount
	}

	return this.Current
}
