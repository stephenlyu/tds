package periodmanager

import (
	"testing"
	. "github.com/stephenlyu/tds/period"
	"fmt"
)


var BASIC_PERIODS = []Period {PERIOD_M, PERIOD_M5, PERIOD_D}

func outputPeriods(periods []Period) {
	for i, p := range periods {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s", p.ShortName())
	}
	fmt.Println("")
}

func TestDefaultPeriodManager_GetPeriodDependencies(t *testing.T) {
	pm := NewDefaultPeriodManager(BASIC_PERIODS)
	periods := pm.GetPeriodDependencies(PeriodFromStringUnsafe("M10"))
	outputPeriods(periods)
}

func TestDefaultPeriodManager_GetOrderedPeriods(t *testing.T) {
	pm := NewDefaultPeriodManager(BASIC_PERIODS)
	pm.AddPeriod(PeriodFromStringUnsafe("M10"))
	pm.AddPeriod(PeriodFromStringUnsafe("M10"))
	pm.AddPeriod(PeriodFromStringUnsafe("Y3"))
	pm.AddPeriod(PeriodFromStringUnsafe("W3"))
	pm.AddPeriod(PeriodFromStringUnsafe("D2"))
	pm.AddPeriod(PeriodFromStringUnsafe("M30"))
	pm.AddPeriod(PeriodFromStringUnsafe("D3"))
	pm.AddPeriod(PeriodFromStringUnsafe("M60"))

	periods := pm.GetOrderedPeriods()
	outputPeriods(periods)
}
