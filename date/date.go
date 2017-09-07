package date

import (
	"time"
)

func GetDateDay(date uint64) int {
	d := time.Unix(int64(date) / 1000, (int64(date) % 1000) * int64(time.Millisecond))
	return d.Year() * 10000 + int(d.Month()) * 100 + d.Day()
}

func GetDateWeek(date uint64) int {
	d := time.Unix(int64(date) / 1000, (int64(date) % 1000) * int64(time.Millisecond))

	y, week := d.ISOWeek()

	return y * 100 + week
}

func GetDateMonth(date uint64) int {
	d := time.Unix(int64(date) / 1000, (int64(date) % 1000) * int64(time.Millisecond))
	return d.Year() * 100 + int(d.Month())
}

var monthQuarterMap = map[int]int {
	1: 3, 2: 3, 3: 3,
	4: 6, 5: 6, 6: 6,
	7: 9, 8: 9, 9: 9,
	10: 12, 11: 12, 12: 12,
}
func GetDateQuarter(date uint64) int {
	d := time.Unix(int64(date) / 1000, (int64(date) % 1000) * int64(time.Millisecond))
	return d.Year() * 100 + monthQuarterMap[int(d.Month())]
}

func GetDateYear(date uint64) int {
	d := time.Unix(int64(date) / 1000, (int64(date) % 1000) * int64(time.Millisecond))
	return d.Year()
}
