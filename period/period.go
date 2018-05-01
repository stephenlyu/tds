package period

import (
	"regexp"
	"strings"
	"errors"
	"strconv"
	"fmt"
	"github.com/stephenlyu/tds/util"
)

type PeriodUnit int

const (
	PERIOD_UNIT_MINUTE = iota  	// M, MINUTE
	PERIOD_UNIT_DAY				// D, DAY
	PERIOD_UNIT_WEEK			// W, WEEK
	PERIOD_UNIT_MONTH			// N, MONTH
	PERIOD_UNIT_QUARTER			// Q, QUARTER
	PERIOD_UNIT_YEAR			// Y, YEAR
)

var (
	_, PERIOD_M = PeriodFromString("M1")
	_, PERIOD_M5 = PeriodFromString("M5")
	_, PERIOD_M15 = PeriodFromString("M15")
	_, PERIOD_M30 = PeriodFromString("M30")
	_, PERIOD_M60 = PeriodFromString("M60")
	_, PERIOD_D = PeriodFromString("D1")
	_, PERIOD_W = PeriodFromString("W1")
	_, PERIOD_MONTH = PeriodFromString("N1")
	_, PERIOD_Q = PeriodFromString("Q1")
	_, PERIOD_Y = PeriodFromString("Y1")
)

type Period interface {
	Name() string
	ShortName() string
	DisplayName() string
	Unit() PeriodUnit
	UnitCount() int
	Eq(other Period) bool
	Lt(other Period) bool
	Gt(other Period) bool
	CanConvertTo(other Period) bool
	CanConvertFrom(other Period) bool
	BasicMergePeriod() Period
	KLinePerDay(nMinutes int) int
}

type period struct {
	unit PeriodUnit
	unitCount int

	shortName, name string
}

// periodStr format: M1 or MINUTE1
func PeriodFromString(periodStr string) (error, Period) {
	regExp, _ := regexp.Compile("^([A-Z]+)([0-9]+)$")

	result := regExp.FindSubmatch([]byte(strings.ToUpper(periodStr)))
	if result == nil {
		return errors.New("bad period string"), nil
	}

	unitStr := string(result[1])
	nUnit, _ := strconv.Atoi(string(result[2]))

	if nUnit == 0 {
		return errors.New("bad unit count"), nil
	}

	var p *period

	switch unitStr {
	case "M":
		fallthrough
	case "MINUTE":
		p = &period{unit: PERIOD_UNIT_MINUTE, unitCount: nUnit}
	case "D":
		fallthrough
	case "DAY":
		p = &period{unit: PERIOD_UNIT_DAY, unitCount: nUnit}
	case "W":
		fallthrough
	case "WEEK":
		p = &period{unit: PERIOD_UNIT_WEEK, unitCount: nUnit}
	case "N":
		fallthrough
	case "MONTH":
		p = &period{unit: PERIOD_UNIT_MONTH, unitCount: nUnit}
	case "Q":
		fallthrough
	case "QUARTER":
		p = &period{unit: PERIOD_UNIT_QUARTER, unitCount: nUnit}
	case "Y":
		fallthrough
	case "YEAR":
		p = &period{unit: PERIOD_UNIT_YEAR, unitCount: nUnit}
	default:
		return errors.New("bad period string"), nil
	}

	p.shortName = p.calcShortName()
	p.name = p.calcName()

	return nil, p
}

func PeriodFromStringUnsafe(periodStr string) Period {
	_, ret := PeriodFromString(periodStr)
	return ret
}

func (this *period) calcName() string {
	switch (this.unit) {
	case PERIOD_UNIT_MINUTE:
		return fmt.Sprintf("MINUTE%d", this.unitCount)
	case PERIOD_UNIT_DAY:
		return fmt.Sprintf("DAY%d", this.unitCount)
	case PERIOD_UNIT_WEEK:
		return fmt.Sprintf("WEEK%d", this.unitCount)
	case PERIOD_UNIT_MONTH:
		return fmt.Sprintf("MONTH%d", this.unitCount)
	case PERIOD_UNIT_QUARTER:
		return fmt.Sprintf("QUARTER%d", this.unitCount)
	case PERIOD_UNIT_YEAR:
		return fmt.Sprintf("YEAR%d", this.unitCount)
	}
	return ""
}

func (this *period) calcShortName() string {
	switch (this.unit) {
	case PERIOD_UNIT_MINUTE:
		return fmt.Sprintf("M%d", this.unitCount)
	case PERIOD_UNIT_DAY:
		return fmt.Sprintf("D%d", this.unitCount)
	case PERIOD_UNIT_WEEK:
		return fmt.Sprintf("W%d", this.unitCount)
	case PERIOD_UNIT_MONTH:
		return fmt.Sprintf("N%d", this.unitCount)
	case PERIOD_UNIT_QUARTER:
		return fmt.Sprintf("Q%d", this.unitCount)
	case PERIOD_UNIT_YEAR:
		return fmt.Sprintf("Y%d", this.unitCount)
	}
	return ""
}

func (this *period) Name() string {
	return this.name
}

func (this *period) ShortName() string {
	return this.shortName
}

func (this *period) DisplayName() string {
	switch (this.unit) {
	case PERIOD_UNIT_MINUTE:
		return fmt.Sprintf("%d分钟", this.unitCount)
	case PERIOD_UNIT_DAY:
		if this.unitCount == 1 {
			return "日线"
		}
		return fmt.Sprintf("%d日线", this.unitCount)
	case PERIOD_UNIT_WEEK:
		if this.unitCount == 1 {
			return "周线"
		}
		return fmt.Sprintf("%d周线", this.unitCount)
	case PERIOD_UNIT_MONTH:
		if this.unitCount == 1 {
			return "月线"
		}
		return fmt.Sprintf("%d月线", this.unitCount)
	case PERIOD_UNIT_QUARTER:
		if this.unitCount == 1 {
			return "季线"
		}
		return fmt.Sprintf("%d季线", this.unitCount)
	case PERIOD_UNIT_YEAR:
		if this.unitCount == 1 {
			return "年线"
		}
		return fmt.Sprintf("%d年线", this.unitCount)
	}
	return ""
}

func (this *period) Unit() PeriodUnit {
	return this.unit
}

func (this *period) UnitCount() int {
	return this.unitCount
}

func (this *period) Eq(other Period) bool {
	return this.Unit() == other.Unit() && this.UnitCount() == other.UnitCount()
}

func (this *period) Lt(other Period) bool {
	return this.Unit() < other.Unit() || this.Unit() == other.Unit() && this.UnitCount() < other.UnitCount()
}

func (this *period) Gt(other Period) bool {
	return this.Unit() > other.Unit() || this.Unit() == other.Unit() && this.UnitCount() > other.UnitCount()
}

func (this *period) CanConvertTo(other Period) bool {
	switch other.Unit() {
	case PERIOD_UNIT_MINUTE:
		switch this.Unit() {
		case PERIOD_UNIT_MINUTE:
			return other.UnitCount() % this.UnitCount() == 0
		default:
			return false
		}
	case PERIOD_UNIT_DAY:
		switch this.Unit() {
		case PERIOD_UNIT_MINUTE:
			return other.UnitCount() == 1
		case PERIOD_UNIT_DAY:
			return other.UnitCount() % this.UnitCount() == 0
		default:
			return false
		}
	case PERIOD_UNIT_WEEK:
		switch this.Unit() {
		case PERIOD_UNIT_DAY:
			return this.UnitCount() == 1 && other.UnitCount() == 1
		case PERIOD_UNIT_WEEK:
			return other.UnitCount() % this.UnitCount() == 0
		default:
			return false
		}
	case PERIOD_UNIT_MONTH:
		switch this.Unit() {
		case PERIOD_UNIT_DAY:
			return this.UnitCount() == 1 && other.UnitCount() == 1
		case PERIOD_UNIT_MONTH:
			return other.UnitCount() % this.UnitCount() == 0
		default:
			return false
		}
	case PERIOD_UNIT_QUARTER:
		switch this.Unit() {
		case PERIOD_UNIT_DAY:
			fallthrough
		case PERIOD_UNIT_MONTH:
			return this.UnitCount() == 1 && other.UnitCount() == 1
		case PERIOD_UNIT_QUARTER:
			return other.UnitCount() % this.UnitCount() == 0
		default:
			return false
		}
	case PERIOD_UNIT_YEAR:
		switch this.Unit() {
		case PERIOD_UNIT_DAY:
			fallthrough
		case PERIOD_UNIT_MONTH:
			fallthrough
		case PERIOD_UNIT_QUARTER:
			return this.UnitCount() == 1 && other.UnitCount() == 1
		case PERIOD_UNIT_YEAR:
			return other.UnitCount() % this.UnitCount() == 0
		default:
			return false
		}
	}

	return true
}

func (this *period) CanConvertFrom(other Period) bool {
	return other.CanConvertTo(this)
}

func (this *period) BasicMergePeriod() Period {
	switch this.Unit() {
	case PERIOD_UNIT_MINUTE:
		return PERIOD_M
	case PERIOD_UNIT_DAY:
		switch this.UnitCount() {
		case 1:
			return PERIOD_M
		default:
			return PERIOD_D
		}
	case PERIOD_UNIT_WEEK:
		switch this.UnitCount() {
		case 1:
			return PERIOD_D
		default:
			return PERIOD_W
		}
	case PERIOD_UNIT_MONTH:
		switch this.UnitCount() {
		case 1:
			return PERIOD_D
		default:
			return PERIOD_MONTH
		}
	case PERIOD_UNIT_QUARTER:
		switch this.UnitCount() {
		case 1:
			return PERIOD_MONTH
		default:
			return PERIOD_Q
		}
	case PERIOD_UNIT_YEAR:
		switch this.UnitCount() {
		case 1:
			return PERIOD_D
		default:
			return PERIOD_Y
		}
	}

	util.Assert(false, "Unreachable code")
	return nil
}

func (this *period) KLinePerDay(nMinutes int) int {
	if nMinutes == 0 {
		nMinutes = 240
	}
	switch this.Unit() {
	case PERIOD_UNIT_MINUTE:
		return (nMinutes + this.unitCount - 1) / this.unitCount
	default:
		return 1
	}
}
