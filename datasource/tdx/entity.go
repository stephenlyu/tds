package tdxdatasource

import (
	"encoding/binary"
	"errors"
	"math"
	"unsafe"
	"github.com/stephenlyu/tds/entity"
	"fmt"
	"strings"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/util"
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

func TDXRecordFromBytes(p period.Period, data []byte, r *TDXRecord) error {
	if len(data) != TDX_RECORD_SIZSE {
		return errors.New("less record bytes")
	}

	if r == nil {
		return errors.New("bad record argument")
	}

	if p.Unit() == period.PERIOD_UNIT_MINUTE {
		r.Date = binary.LittleEndian.Uint32(data[0:4])
		r.Open = int32(math.Float32frombits(binary.LittleEndian.Uint32(data[4:8])) * 1000)
		r.High = int32(math.Float32frombits(binary.LittleEndian.Uint32(data[8:12])) * 1000)
		r.Low = int32(math.Float32frombits(binary.LittleEndian.Uint32(data[12:16])) * 1000)
		r.Close = int32(math.Float32frombits(binary.LittleEndian.Uint32(data[16:20])) * 1000)
		r.Amount = math.Float32frombits(binary.LittleEndian.Uint32(data[20:24])) * 100
		r.Volume = float32(binary.LittleEndian.Uint32(data[24:28])) * 100
	} else if p.Unit() == period.PERIOD_UNIT_DAY {
		r.Date = binary.LittleEndian.Uint32(data[0:4])
		r.Open = int32(binary.LittleEndian.Uint32(data[4:8])) * 10
		r.High = int32(binary.LittleEndian.Uint32(data[8:12])) * 10
		r.Low = int32(binary.LittleEndian.Uint32(data[12:16])) * 10
		r.Close = int32(binary.LittleEndian.Uint32(data[16:20])) * 10
		r.Amount = math.Float32frombits(binary.LittleEndian.Uint32(data[20:24]))
		r.Volume = float32(binary.LittleEndian.Uint32(data[24:28]))
	} else {
		util.Assert(false, "Unsupported period")
	}
	return nil
}

func (this *TDXRecord) Bytes(p period.Period) []byte {
	ret := make([]byte, TDX_RECORD_SIZSE)

	if p.Unit() == period.PERIOD_UNIT_MINUTE {
		binary.LittleEndian.PutUint32(ret[0:4], this.Date)
		binary.LittleEndian.PutUint32(ret[4:8], math.Float32bits(float32(this.Open) / 1000))
		binary.LittleEndian.PutUint32(ret[8:12], math.Float32bits(float32(this.High) / 1000))
		binary.LittleEndian.PutUint32(ret[12:16], math.Float32bits(float32(this.Low) / 1000))
		binary.LittleEndian.PutUint32(ret[16:20], math.Float32bits(float32(this.Close) / 1000))
		binary.LittleEndian.PutUint32(ret[20:24], math.Float32bits(this.Amount))
		binary.LittleEndian.PutUint32(ret[24:28], uint32(this.Volume / 100))
	} else if p.Unit() == period.PERIOD_UNIT_DAY {
		binary.LittleEndian.PutUint32(ret[0:4], this.Date)
		binary.LittleEndian.PutUint32(ret[4:8], uint32(this.Open / 10))
		binary.LittleEndian.PutUint32(ret[8:12], uint32(this.High / 10))
		binary.LittleEndian.PutUint32(ret[12:16], uint32(this.Low / 10))
		binary.LittleEndian.PutUint32(ret[16:20], uint32(this.Close / 10))
		binary.LittleEndian.PutUint32(ret[20:24], math.Float32bits(this.Amount))
		binary.LittleEndian.PutUint32(ret[24:28], uint32(this.Volume))
	}
	return ret
}

func SecurityToString(security *entity.Security) string {
	return strings.ToLower(fmt.Sprintf("%s%s", security.Exchange, security.Code))
}
