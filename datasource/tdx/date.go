package tdxdatasource

import (
	"time"
	. "github.com/stephenlyu/tds/period"
)

var local, _ = time.LoadLocation("Asia/Shanghai")

func MinuteDateToTimestamp(v uint32) uint64 {
	dayValue := uint16(v & 0xFFFF)
	minuteValue := uint16((v >> 16) & 0xFFFF)

	year := int((dayValue / 2048) + 2004)
	month := int((dayValue % 2048) / 100)
	day := int((dayValue % 2048) % 100)

	hour := int(minuteValue / 60)
	minute := int(minuteValue % 60)

	date := time.Date(year, time.Month(month), day, hour, minute, 0, 0, local)

	return uint64(date.UnixNano() / int64(time.Millisecond))
}

func DayDateToTimestamp(v uint32) uint64 {
	year := int(v / 10000)
	month := int((v % 10000) / 100)
	day := int((v % 10000) % 100)

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, local)

	return uint64(date.UnixNano() / int64(time.Millisecond))
}

func DateToTimestamp(period Period, v uint32) uint64 {
	if period.Unit() == PERIOD_UNIT_MINUTE {
		return MinuteDateToTimestamp(v)
	}
	return DayDateToTimestamp(v)
}

func TimestampToMinuteDate(ts uint64) uint32 {
	date := time.Unix(int64(ts) / 1000, (int64(ts) % 1000) * int64(time.Millisecond))

	dayValue := uint32(date.Year() - 2004) * 2048 + uint32(date.Month()) * 100 + uint32(date.Day())
	minuteValue := uint32(date.Hour()) * 60 + uint32(date.Minute())
	return minuteValue << 16 | dayValue
}

func TimestampToDayDate(ts uint64) uint32 {
	date := time.Unix(int64(ts) / 1000, (int64(ts) % 1000) * int64(time.Millisecond))
	return uint32(date.Year()) * 10000 + uint32(date.Month()) * 100 + uint32(date.Day())
}

func TimestampToDate(period Period, ts uint64) uint32 {
	if period.Unit() == PERIOD_UNIT_MINUTE {
		return TimestampToMinuteDate(ts)
	}
	return TimestampToDayDate(ts)
}
