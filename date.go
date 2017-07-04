package tds

import (
	"math"
	"fmt"
	"strconv"
	"regexp"
	"errors"
	"time"
)

type Date interface {
	DayString() string
	MinuteString() string
	MinuteDate() uint32
	DayDate() uint32
	PeriodValue(period Period) uint32

	Year() uint32
	Month() uint32
	Day() uint32
	Hour() uint32
	Minute() uint32
	Second() uint32

	Eq(other Date) bool
	Lt(other Date) bool
	Gt(other Date) bool
}

type date struct {
	year uint16
	month uint8
	day uint8
	hour uint8
	minute uint8
	second uint8
}

const MAX_DATE = math.MaxUint32

func (this *date) DayString() string {
	return fmt.Sprintf("%04d%02d%02d", this.year, this.month, this.day)
}

func (this *date) MinuteString() string {
	return fmt.Sprintf("%04d%02d%02d %02d:%02d:00", this.year, this.month, this.day, this.hour, this.minute)
}

func (this *date) MinuteDate() uint32 {
	dayValue := uint32(this.year - 2004) * 2048 + uint32(this.month) * 100 + uint32(this.day)
	minuteValue := uint32(this.hour) * 60 + uint32(this.minute)
	return minuteValue << 16 | dayValue
}

func (this *date) DayDate() uint32 {
	return uint32(this.year) * 10000 + uint32(this.month) * 100 + uint32(this.day)
}

func (this *date) PeriodValue(period Period) uint32 {
	if period.Unit() == PERIOD_UNIT_MINUTE {
		return this.MinuteDate()
	}
	return this.DayDate()
}

func (this *date) Year() uint32 {return uint32(this.year)}
func (this *date) Month() uint32 {return uint32(this.month)}
func (this *date) Day() uint32 {return uint32(this.day)}
func (this *date) Hour() uint32 {return uint32(this.hour)}
func (this *date) Minute() uint32 {return uint32(this.minute)}
func (this *date) Second() uint32 {return uint32(this.second)}

func (this *date) Eq(other Date) bool {
	return this.Year() == other.Year() && this.Month() == other.Month() && this.Day() == other.Day() &&
	this.Hour() == other.Hour() && this.Minute() == other.Minute() && this.Second() == other.Second()
}

func (this *date) Lt(other Date) bool {
	return this.Year() < other.Year() || this.Month() < other.Month() || this.Day() < other.Day() ||
		this.Hour() < other.Hour() || this.Minute() < other.Minute() || this.Second() < other.Second()
}

func (this *date) Gt(other Date) bool {
	return this.Year() > other.Year() || this.Month() > other.Month() || this.Day() > other.Day() ||
		this.Hour() > other.Hour() || this.Minute() > other.Minute() || this.Second() > other.Second()
}

func NewDateFromMinuteDate(v uint32) Date {
	dayValue := uint16(v & 0xFFFF)
	minuteValue := uint16((v >> 16) & 0xFFFF)

	year := (dayValue / 2048) + 2004
	month := (dayValue % 2048) / 100
	day := (dayValue % 2048) % 100

	hour := minuteValue / 60
	minute := minuteValue % 60
	return &date{uint16(year), uint8(month), uint8(day), uint8(hour), uint8(minute), 0}
}

func NewDateFromDayDate(v uint32) Date {
	year := v / 10000
	month := (v % 10000) / 100
	day := (v % 10000) % 100
	return &date{uint16(year), uint8(month), uint8(day), 0, 0, 0}
}

func NewPeriodDate(period Period, v uint32) Date {
	if period.Unit() == PERIOD_UNIT_MINUTE {
		return NewDateFromMinuteDate(v)
	}
	return NewDateFromDayDate(v)
}

func NewDateFromDayString(s string) (error, Date) {
	ret, err := strconv.ParseUint(s, 10, 64)
	return err, NewDateFromDayDate(uint32(ret))
}

func NewDateFromMinuteString(s string) (error, Date) {
	regExp, err := regexp.Compile("^([0-9]{4})([0-9]{2})([0-9]{2}) ([0-9]{2}):([0-9]{2}):([0-9]{2})$")
	if err != nil {
		return errors.New("bad regexp pattern"), nil
	}

	result := regExp.FindSubmatch([]byte(s))
	if result == nil {
		return errors.New("bad minute string"), nil
	}

	year, _ := strconv.Atoi(string(result[1]))
	month, _ := strconv.Atoi(string(result[2]))
	day, _ := strconv.Atoi(string(result[3]))
	hour, _ := strconv.Atoi(string(result[4]))
	minute, _ := strconv.Atoi(string(result[5]))
	second, _ := strconv.Atoi(string(result[6]))

	if year < 2004 || year > 2004 + 511 {
		return errors.New("bad year"), nil
	}

	if month <= 0 || month > 12 {
		return errors.New("bad month"), nil
	}

	if day <= 0 || day > 31 {
		return errors.New("bad day"), nil
	}

	if hour < 0 || hour >= 24 {
		return errors.New("bad hour"), nil
	}

	if minute < 0 || minute > 59 {
		return errors.New("bad minute"), nil
	}

	if second < 0 || second > 59 {
		return errors.New("bad second"), nil
	}

	return nil, &date{uint16(year), uint8(month), uint8(day), uint8(hour), uint8(minute), uint8(second)}
}

func GetDateDay(period Period, date uint32) uint32 {
	if period.Unit() == PERIOD_UNIT_MINUTE {
		return NewDateFromMinuteDate(date).MinuteDate()
	}
	return date
}

func GetDateWeek(period Period, dateValue uint32) uint32 {
	var date Date
	if period.Unit() == PERIOD_UNIT_MINUTE {
		date = NewDateFromMinuteDate(dateValue)
	} else {
		date = NewDateFromDayDate(dateValue)
	}

	d := time.Date(int(date.Year()), time.Month(date.Month()), int(date.Day()), 0, 0, 0, 0, time.UTC)

	y, week := d.ISOWeek()

	return uint32(y * 100 + week)
}

func GetDateMonth(period Period, dateValue uint32) uint32 {
	if period.Unit() == PERIOD_UNIT_MINUTE {
		date := NewDateFromMinuteDate(dateValue)
		return date.Year() * 100 + date.Month()
	}
	return dateValue / 100
}

var monthQuarterMap = map[int]uint32 {
	1: 3, 2: 3, 3: 3,
	4: 6, 5: 6, 6: 6,
	7: 9, 8: 9, 9: 9,
	10: 12, 11: 12, 12: 12,
}
func GetDateQuarter(period Period, dateValue uint32) uint32 {
	var date Date
	if period.Unit() == PERIOD_UNIT_MINUTE {
		date = NewDateFromMinuteDate(dateValue)
	} else {
		date = NewDateFromDayDate(dateValue)
	}

	return date.Year() * 100 + monthQuarterMap[int(date.Month())]
}

func GetDateYear(period Period, dateValue uint32) uint32 {
	var date Date
	if period.Unit() == PERIOD_UNIT_MINUTE {
		date = NewDateFromMinuteDate(dateValue)
	} else {
		date = NewDateFromDayDate(dateValue)
	}

	return date.Year()
}
