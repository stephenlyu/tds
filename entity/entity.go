package entity

import (
	"unsafe"
	"github.com/stephenlyu/tds/date"
	"fmt"
	"github.com/golang/protobuf/proto"
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

func (this *Record) ToProtoBytes() ([]byte, error) {
	pr := ProtoRecord{
		Date: int64(this.Date),
		Open: this.Open,
		Close: this.Close,
		High: this.High,
		Low: this.Low,
		Volume: this.Volume,
		Amount: this.Amount,
	}

	return proto.Marshal(&pr)
}

func RecordFromProtoBytes(bytes []byte) (*Record, error) {
	var pr ProtoRecord
	err := proto.Unmarshal(bytes, &pr)
	if err != nil {
		return nil, err
	}

	return &Record{
		Date: uint64(pr.GetDate()),
		Open: pr.GetOpen(),
		Close: pr.GetClose(),
		High: pr.GetHigh(),
		Low: pr.GetLow(),
		Volume: pr.GetVolume(),
		Amount: pr.GetAmount(),
	}, nil
}

type Tick struct {
	Code string
	Timestamp uint64
	HighLimited float64
	LowLimited float64
	Price float64
	Position float64
	Settle float64
	Open float64
	Close float64
	High float64
	Low float64
	Volume float64
	TotalVolume float64
	Amount float64
	TotalAmount float64
	PreSettle float64
	PrePosition float64
	PreClose float64
	AskPrices []float64
	AskVolumes []float64
	BidPrices []float64
	BidVolumes []float64
}

func (this *Tick) GetDate() string {
	return date.Timestamp2SecondString(this.Timestamp)
}

func (this *Tick) ToProtoBytes() ([]byte, error) {
	pr := ProtoTick{
		Code: this.Code,
		Timestamp: int64(this.Timestamp),
		HighLimited: this.HighLimited,
		LowLimited: this.LowLimited,
		Price: this.Price,
		Position: this.Position,
		Settle: this.Settle,
		Open: this.Open,
		Close: this.Close,
		High: this.High,
		Low: this.Low,
		Volume: this.Volume,
		TotalVolume: this.TotalVolume,
		Amount: this.Amount,
		TotalAmount: this.TotalAmount,
		PreSettle: this.PreSettle,
		PrePosition: this.PrePosition,
		PreClose: this.PreClose,
		AskPrices: this.AskPrices,		// 不拷贝
		AskVolumes: this.AskVolumes,
		BidPrices: this.BidPrices,
		BidVolumes: this.BidVolumes,
	}

	return proto.Marshal(&pr)
}

func TickFromProtoBytes(bytes []byte) (*Tick, error) {
	var pr ProtoTick
	err := proto.Unmarshal(bytes, &pr)
	if err != nil {
		return nil, err
	}

	return &Tick{
		Code: pr.Code,
		Timestamp: uint64(pr.Timestamp),
		HighLimited: pr.HighLimited,
		LowLimited: pr.LowLimited,
		Price: pr.Price,
		Position: pr.Position,
		Settle: pr.Settle,
		Open: pr.Open,
		Close: pr.Close,
		High: pr.High,
		Low: pr.Low,
		Volume: pr.Volume,
		TotalVolume: pr.TotalVolume,
		Amount: pr.Amount,
		TotalAmount: pr.TotalAmount,
		PreSettle: pr.PreSettle,
		PrePosition: pr.PrePosition,
		PreClose: pr.PreClose,
		AskPrices: pr.AskPrices,		// 不拷贝
		AskVolumes: pr.AskVolumes,
		BidPrices: pr.BidPrices,
		BidVolumes: pr.BidVolumes,
	}, nil
}
