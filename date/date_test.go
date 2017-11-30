package date_test

import (
	"testing"
	"github.com/stephenlyu/tds/util"
	"fmt"
	"github.com/stephenlyu/tds/date"
)

func TestDate(t *testing.T) {
	now := util.Tick()
	fmt.Println(date.GetDateDay(now))
	fmt.Println(date.GetDateMonth(now))
	fmt.Println(date.GetDateWeek(now))
	fmt.Println(date.GetDateYear(now))
	fmt.Println(date.GetDateQuarter(now))
	fmt.Println(date.GetNowString())
	fmt.Println(date.GetTodayString())
}