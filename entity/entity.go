package entity

import (
	"encoding/binary"
	"math"
	"errors"
	"unsafe"
	"github.com/stephenlyu/tds/date"
	"fmt"
)


// After price adjusted, it maybe a negative value
type Record struct {
	Date uint64			// UTC毫秒数
	Open int32			// 开盘价，精确到分
	Close int32
	High int32
	Low int32
	Volume float32
	Amount float32
}

type InfoExItem struct {
	Date uint32					`json:"date"`
	Bonus float32				`json:"bonus"`
	DeliveredShares float32		`json:"delivered_shares"`
	RationedSharePrice float32	`json:"rationed_share_price"`
	RationedShares float32		`json:"rationed_shares"`
}

const recordSize = int(unsafe.Sizeof(Record{}))

func (this *Record) Eq(that *Record) bool {
	return this.Date == that.Date && this.Open == that.Open && this.Close == that.Close && this.High == that.High &&
	this.Low == that.Low && this.Volume == that.Volume && this.Amount == that.Amount
}

func (this *Record) GetDateString() string {
	return date.Timestamp2SecondString(this.Date)
}

func (this *Record) GetOpen() float32 {
	return float32(this.Open) / 1000.0
}

func (this *Record) GetClose() float32 {
	return float32(this.Close) / 1000.0
}

func (this *Record) GetLow() float32 {
	return float32(this.Low) / 1000.0
}

func (this *Record) GetHigh() float32 {
	return float32(this.High) / 1000.0
}

func (this *Record) GetAmount() float32 {
	return this.Amount
}

func (this *Record) GetVolume() float32 {
	return this.Volume
}

func (this *Record) String() string {
	return fmt.Sprintf(`Record {Date: %s Open: %.02f Close: %.02f Low: %.02f High: %.02f Amount: %.02f Volume: %.02f}`, this.GetDateString(), this.GetOpen(), this.GetClose(), this.GetLow(), this.GetHigh(), this.GetAmount(), this.GetVolume())
}

func RecordFromBytes(data []byte, r *Record) error {
	if len(data) != recordSize {
		return errors.New("less record bytes")
	}

	if r == nil {
		return errors.New("bad record argument")
	}

	r.Date = binary.LittleEndian.Uint64(data[0:8])
	r.Open = int32(binary.LittleEndian.Uint32(data[8:12]))
	r.Close = int32(binary.LittleEndian.Uint32(data[12:16]))
	r.High = int32(binary.LittleEndian.Uint32(data[16:20]))
	r.Low = int32(binary.LittleEndian.Uint32(data[20:24]))
	r.Volume = math.Float32frombits(binary.LittleEndian.Uint32(data[24:28]))
	r.Amount = math.Float32frombits(binary.LittleEndian.Uint32(data[28:32]))
	return nil
}

func (this *Record) Bytes() []byte {
	ret := make([]byte, recordSize)

	binary.LittleEndian.PutUint64(ret[0:8], this.Date)
	binary.LittleEndian.PutUint32(ret[8:12], uint32(this.Open))
	binary.LittleEndian.PutUint32(ret[12:16], uint32(this.Close))
	binary.LittleEndian.PutUint32(ret[16:20], uint32(this.High))
	binary.LittleEndian.PutUint32(ret[20:24], uint32(this.Low))
	binary.LittleEndian.PutUint32(ret[24:28], math.Float32bits(this.Volume))
	binary.LittleEndian.PutUint32(ret[28:32], math.Float32bits(this.Amount))
	return ret
}
