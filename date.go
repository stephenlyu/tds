package tds

import (
	"math"
	"fmt"
	"strconv"
	"regexp"
	"errors"
)

type Date uint32

const MAX_DATE = math.MaxUint32

func (this Date) DayString() string {
	return fmt.Sprintf("%d", this)
}

func (this Date) MinuteString() string {
	dayValue := uint16(this & 0xFFFF)
	minuteValue := uint16((this >> 16) & 0xFFFF)

	year := (dayValue / 2048) + 2004
	month := (dayValue % 2048) / 100
	day := (dayValue % 2048) % 100

	hour := minuteValue / 60
	minute := minuteValue % 60

	return fmt.Sprintf("%04d%02d%02d %02d:%02d:00", year, month, day, hour, minute)
}

func FromDayString(s string) (error, Date) {
	ret, err := strconv.ParseUint(s, 10, 64)
	return err, Date(ret)
}

func FromMinuteString(s string) (error, Date) {
	regExp, err := regexp.Compile("^([0-9]{4})([0-9]{2})([0-9]{2}) ([0-9]{2}):([0-9]{2}):([0-9]{2})$")
	if err != nil {
		return errors.New("bad minute string"), 0
	}

	result := regExp.FindSubmatch([]byte(s))
	if result != nil {
		return errors.New("bad minute string"), 0
	}

	year, _ := strconv.Atoi(string(result[1]))
	month, _ := strconv.Atoi(string(result[2]))
	day, _ := strconv.Atoi(string(result[3]))
	hour, _ := strconv.Atoi(string(result[4]))
	minute, _ := strconv.Atoi(string(result[5]))
	second, _ := strconv.Atoi(string(result[6]))

	if year < 2004 || year > 2004 + 511 {
		return errors.New("bad year"), 0
	}

	if month <= 0 || month > 12 {
		return errors.New("bad month"), 0
	}

	if day <= 0 || day > 31 {
		return errors.New("bad day"), 0
	}

	if hour <= 0 || year >= 24 {
		return errors.New("bad hour"), 0
	}

	if minute <= 0 || minute > 59 {
		return errors.New("bad minute"), 0
	}

	if second <= 0 || second > 59 {
		return errors.New("bad second"), 0
	}

	dayValue := year * 2048 + month * 100 + day
	minuteValue := hour * 60 + minute

	return nil, Date((minuteValue << 16) | dayValue)
}
