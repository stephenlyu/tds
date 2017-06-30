package tds

import (
	"regexp"
	"strings"
	"errors"
	"strconv"
	"fmt"
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

type Period interface {
	Name() string
	ShortName() string
}

type period struct {
	Unit PeriodUnit
	UnitCount int
}

// periodStr format: M1 or MINUTE1
func FromString(periodStr string) (error, Period) {
	regExp, _ := regexp.Compile("^([A-Z]+)([0-9]+)$")

	result := regExp.FindSubmatch([]byte(strings.ToUpper(periodStr)))
	if result != nil {
		return errors.New("bad period string"), nil
	}

	unitStr := string(result[1])
	nUnit, _ := strconv.Atoi(string(result[2]))

	if nUnit == 0 {
		return errors.New("bad period string"), nil
	}

	switch unitStr {
	case "M":
		fallthrough
	case "MINUTE":
		return nil, &period{PERIOD_UNIT_MINUTE, nUnit}
	case "D":
		fallthrough
	case "DAY":
		return nil, &period{PERIOD_UNIT_DAY, nUnit}
	case "W":
		fallthrough
	case "WEEK":
		return nil, &period{PERIOD_UNIT_WEEK, nUnit}
	case "N":
		fallthrough
	case "MONTH":
		return nil, &period{PERIOD_UNIT_MONTH, nUnit}
	case "Q":
		fallthrough
	case "QUARTER":
		return nil, &period{PERIOD_UNIT_QUARTER, nUnit}
	case "Y":
		fallthrough
	case "YEAR":
		return nil, &period{PERIOD_UNIT_YEAR, nUnit}
	}

	return errors.New("bad period string"), nil
}

func (this *period) Name() string {
	switch (this.Unit) {
	case PERIOD_UNIT_MINUTE:
		return fmt.Sprintf("MINUTE%d", this.UnitCount)
	case PERIOD_UNIT_DAY:
		return fmt.Sprintf("DAY%d", this.UnitCount)
	case PERIOD_UNIT_WEEK:
		return fmt.Sprintf("WEEK%d", this.UnitCount)
	case PERIOD_UNIT_MONTH:
		return fmt.Sprintf("MONTH%d", this.UnitCount)
	case PERIOD_UNIT_QUARTER:
		return fmt.Sprintf("QUARTER%d", this.UnitCount)
	case PERIOD_UNIT_YEAR:
		return fmt.Sprintf("YEAR%d", this.UnitCount)
	}
	return ""
}

func (this *period) ShortName() string {
	switch (this.Unit) {
	case PERIOD_UNIT_MINUTE:
		return fmt.Sprintf("M%d", this.UnitCount)
	case PERIOD_UNIT_DAY:
		return fmt.Sprintf("D%d", this.UnitCount)
	case PERIOD_UNIT_WEEK:
		return fmt.Sprintf("W%d", this.UnitCount)
	case PERIOD_UNIT_MONTH:
		return fmt.Sprintf("N%d", this.UnitCount)
	case PERIOD_UNIT_QUARTER:
		return fmt.Sprintf("Q%d", this.UnitCount)
	case PERIOD_UNIT_YEAR:
		return fmt.Sprintf("Y%d", this.UnitCount)
	}
	return ""
}
