package entity

import (
	"unsafe"
	"github.com/stephenlyu/tds/date"
	"fmt"
)


// After price adjusted, it maybe a negative value
type Record struct {
	Date uint64			`bson:"_id"`	// UTC毫秒数
	Open float64						// 开盘价，精确到分
	Close float64
	High float64
	Low float64
	Volume float64
	Amount float64
}

type InfoExItem struct {
	Date uint32					`json:"date"`
	Bonus float64				`json:"bonus"`
	DeliveredShares float64		`json:"delivered_shares"`
	RationedSharePrice float64	`json:"rationed_share_price"`
	RationedShares float64		`json:"rationed_shares"`
}

const recordSize = int(unsafe.Sizeof(Record{}))

func (this *Record) Eq(that *Record) bool {
	return this.Date == that.Date && this.Open == that.Open && this.Close == that.Close && this.High == that.High &&
	this.Low == that.Low && this.Volume == that.Volume && this.Amount == that.Amount
}

func (this *Record) GetUTCDate() uint64 {
	return this.Date
}

func (this *Record) GetDate() string {
	return date.Timestamp2SecondString(this.Date)
}

func (this *Record) GetOpen() float64 {
	return this.Open
}

func (this *Record) GetClose() float64 {
	return this.Close
}

func (this *Record) GetLow() float64 {
	return this.Low
}

func (this *Record) GetHigh() float64 {
	return this.High
}

func (this *Record) GetAmount() float64 {
	return this.Amount
}

func (this *Record) GetVolume() float64 {
	return this.Volume
}

func (this *Record) String() string {
	return fmt.Sprintf(`Record {Date: %s Open: %.02f Close: %.02f Low: %.02f High: %.02f Amount: %.02f Volume: %.02f}`, this.GetDate(), this.GetOpen(), this.GetClose(), this.GetLow(), this.GetHigh(), this.GetAmount(), this.GetVolume())
}
