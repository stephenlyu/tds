package tdxdatasource

import (
	"github.com/stephenlyu/tds/datasource"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/util"
)

const MINUTE = 60 * 1000
const MINUTE_1300 = 13 * 60 * MINUTE
const MINUTE_1130 = (11 * 60 + 30) * MINUTE
const DAY = 24 * 60 * MINUTE

type tdxMarshaller struct {
	period period.Period
}

func NewMarshaller(period period.Period) datasource.RecordMarshaller {
	return &tdxMarshaller{period}
}

func (this *tdxMarshaller) ToBytes(record *entity.Record) ([]byte, error) {
	date := record.Date
	if this.period.Unit() == period.PERIOD_UNIT_MINUTE {
		date += MINUTE
		if date % DAY == MINUTE_1130 {
			date += 90 * MINUTE
		}
	}

	tRecord := TDXRecord{
		Date: TimestampToDate(this.period, date),
		Open: int32(util.Round(record.Open * 1000, 0)),
		Close: int32(util.Round(record.Close * 1000, 0)),
		High: int32(util.Round(record.High * 1000, 0)),
		Low: int32(util.Round(record.Low * 1000, 0)),
		Amount: float32(record.Amount),
		Volume: float32(record.Volume),
	}
	return tRecord.Bytes(this.period), nil
}

func (this *tdxMarshaller) FromBytes(bytes []byte, record *entity.Record) error {
	tRecord := TDXRecord{}
	err := TDXRecordFromBytes(this.period, bytes, &tRecord)
	if err != nil {
		return err
	}

	record.Date = DateToTimestamp(this.period, tRecord.Date)
	record.Open = float64(tRecord.Open) / 1000
	record.Close = float64(tRecord.Close) / 1000
	record.High = float64(tRecord.High) / 1000
	record.Low = float64(tRecord.Low) / 1000
	record.Amount = float64(tRecord.Amount)
	record.Volume = float64(tRecord.Volume)

	if this.period.Unit() == period.PERIOD_UNIT_MINUTE {
		if record.Date % DAY == MINUTE_1300 {
			record.Date -= 90 * MINUTE
		}
		record.Date -= MINUTE
	}

	return nil
}
