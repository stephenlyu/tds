package tdxdatasource

import (
	"github.com/stephenlyu/tds/datasource"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/period"
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
		Open: record.Open,
		Close: record.Close,
		High: record.High,
		Low: record.Low,
		Amount: record.Amount,
		Volume: record.Volume,
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
	record.Open = tRecord.Open
	record.Close = tRecord.Close
	record.High = tRecord.High
	record.Low = tRecord.Low
	record.Amount = tRecord.Amount
	record.Volume = tRecord.Volume

	if this.period.Unit() == period.PERIOD_UNIT_MINUTE {
		if record.Date % DAY == MINUTE_1300 {
			record.Date -= 90 * MINUTE
		}
		record.Date -= MINUTE
	}

	return nil
}
