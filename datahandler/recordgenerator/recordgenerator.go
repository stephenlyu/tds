package recordgenerator

import (
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/tradedate"
	"math"
	"github.com/stephenlyu/tds/util"
)

//
// 从TickItem合成1分钟K线程序
//

type RecordGenerator struct {
	Security *entity.Security

	Current *entity.Record		// 当前的K线
}

func NewRecordGenerator(security *entity.Security) *RecordGenerator {
	return &RecordGenerator{
		Security: security,
	}
}

func (this *RecordGenerator) Feed(tick *entity.TickItem) *entity.Record {
	util.Assert(tick.Code == this.Security.String(), "")
	if tick.Price == 0 || tick.Volume == 0 {
		return nil
	}

	ticker := tradedate.ToTradeTicker(this.Security, tick.Timestamp)

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
