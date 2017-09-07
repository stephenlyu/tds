package entity

import (
	"encoding/binary"
	"math"
	"errors"
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

const recordSize = 28

func RecordFromBytes(data []byte, r *Record) error {
	if len(data) != recordSize {
		return errors.New("less record bytes")
	}

	if r == nil {
		return errors.New("bad record argument")
	}

	r.Date = binary.LittleEndian.Uint32(data[0:4])
	r.Open = int32(binary.LittleEndian.Uint32(data[4:8]))
	r.Close = int32(binary.LittleEndian.Uint32(data[8:12]))
	r.High = int32(binary.LittleEndian.Uint32(data[12:16]))
	r.Low = int32(binary.LittleEndian.Uint32(data[16:20]))
	r.Volume = math.Float32frombits(binary.LittleEndian.Uint32(data[20:24]))
	r.Amount = math.Float32frombits(binary.LittleEndian.Uint32(data[24:28]))
	return nil
}

func (this *Record) Bytes() []byte {
	ret := make([]byte, recordSize)

	binary.LittleEndian.PutUint32(ret[0:4], this.Date)
	binary.LittleEndian.PutUint32(ret[4:8], uint32(this.Open))
	binary.LittleEndian.PutUint32(ret[8:12], uint32(this.Close))
	binary.LittleEndian.PutUint32(ret[12:16], uint32(this.High))
	binary.LittleEndian.PutUint32(ret[16:20], uint32(this.Low))
	binary.LittleEndian.PutUint32(ret[20:24], math.Float32bits(this.Volume))
	binary.LittleEndian.PutUint32(ret[24:28], math.Float32bits(this.Amount))
	return ret
}
