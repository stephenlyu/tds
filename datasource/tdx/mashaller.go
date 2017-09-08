package tdxdatasource

import (
	"github.com/stephenlyu/tds/datasource"
	"github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/period"
)

type tdxMarshaller struct {
	period period.Period
}

func NewMarshaller(period period.Period) datasource.RecordMarshaller {
	return &tdxMarshaller{period}
}

func (this *tdxMarshaller) ToBytes(record *entity.Record) ([]byte, error) {
	tRecord := TDXRecord{
		Date: TimestampToDate(this.period, record.Date),
		Open: record.Open,
		Close: record.Close,
		High: record.High,
		Low: record.Low,
		Amount: record.Amount,
		Volume: record.Volume,
	}
	return tRecord.Bytes(), nil
}

func (this *tdxMarshaller) FromBytes(bytes []byte, record *entity.Record) error {
	tRecord := TDXRecord{}
	err := TDXRecordFromBytes(bytes, &tRecord)
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
	return nil
}
