package tdxdatasource

import (
	"encoding/binary"
	"errors"
	"math"
	"unsafe"
	"github.com/stephenlyu/tds/entity"
	"fmt"
	"strings"
)

type TDXRecord struct {
	Date uint32
	Open int32
	Close int32
	High int32
	Low int32
	Volume float32
	Amount float32
	Pad int32
}

const TDX_RECORD_SIZSE = int(unsafe.Sizeof(TDXRecord{}))

func TDXRecordFromBytes(data []byte, r *TDXRecord) error {
	if len(data) != TDX_RECORD_SIZSE {
		return errors.New("less record bytes")
	}

	if r == nil {
		return errors.New("bad record argument")
	}

	r.Date = binary.LittleEndian.Uint32(data[0:4])
	r.Open = int32(math.Float32frombits(binary.LittleEndian.Uint32(data[4:8])) * 1000)
	r.High = int32(math.Float32frombits(binary.LittleEndian.Uint32(data[8:12])) * 1000)
	r.Low = int32(math.Float32frombits(binary.LittleEndian.Uint32(data[12:16])) * 1000)
	r.Close = int32(math.Float32frombits(binary.LittleEndian.Uint32(data[16:20])) * 1000)
	r.Amount = math.Float32frombits(binary.LittleEndian.Uint32(data[20:24])) * 100
	r.Volume = float32(binary.LittleEndian.Uint32(data[24:28])) * 100
	return nil
}

func (this *TDXRecord) Bytes() []byte {
	ret := make([]byte, TDX_RECORD_SIZSE)

	binary.LittleEndian.PutUint32(ret[0:4], this.Date)
	binary.LittleEndian.PutUint32(ret[4:8], math.Float32bits(float32(this.Open) / 1000))
	binary.LittleEndian.PutUint32(ret[8:12], math.Float32bits(float32(this.High) / 1000))
	binary.LittleEndian.PutUint32(ret[12:16], math.Float32bits(float32(this.Low) / 1000))
	binary.LittleEndian.PutUint32(ret[16:20], math.Float32bits(float32(this.Close) / 1000))
	binary.LittleEndian.PutUint32(ret[20:24], math.Float32bits(this.Amount))
	binary.LittleEndian.PutUint32(ret[24:28], uint32(this.Volume / 100))
	return ret
}

func SecurityToString(security *entity.Security) string {
	return strings.ToLower(fmt.Sprintf("%s%s", security.Exchange, security.Code))
}
