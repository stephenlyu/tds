package tdxdatasource

import (
	"encoding/binary"
	"errors"
	"math"
)

type TDXRecord struct {
	Date uint32
	Open int32
	Close int32
	High int32
	Low int32
	Volume float32
	Amount float32
}

const TDX_RECORD_SIZSE = 28

func RecordFromBytes(data []byte, r *TDXRecord) error {
	if len(data) != TDX_RECORD_SIZSE {
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

func (this *TDXRecord) Bytes() []byte {
	ret := make([]byte, TDX_RECORD_SIZSE)

	binary.LittleEndian.PutUint64(ret[0:8], this.Date)
	binary.LittleEndian.PutUint32(ret[8:12], uint32(this.Open))
	binary.LittleEndian.PutUint32(ret[12:16], uint32(this.Close))
	binary.LittleEndian.PutUint32(ret[16:20], uint32(this.High))
	binary.LittleEndian.PutUint32(ret[20:24], uint32(this.Low))
	binary.LittleEndian.PutUint32(ret[24:28], math.Float32bits(this.Volume))
	binary.LittleEndian.PutUint32(ret[28:32], math.Float32bits(this.Amount))
	return ret
}
