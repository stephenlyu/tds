package m1smoother

import (
	"testing"
	"github.com/stephenlyu/tds/entity"
	"fmt"
)

func Test_M1Smoother_Feed(t *testing.T) {
	security := entity.ParseSecurityUnsafe("EOSQFUT.OKEX")

	r0 := entity.Record{Date:1534237140000, Open:4.156, Close:4.159, High: 4.159, Low:4.156, Volume:40, Amount:400}
	r1 := entity.Record{Date:1534237200000, Open:4.156, Close:4.159, High: 4.159, Low:4.156, Volume:40, Amount:400}
	r2 := entity.Record{Date:1534237740000,Open:4.154,Close:4.154,High:4.154,Low:4.154,Volume:996,Amount:9960}
	r3 := entity.Record{Date:1534237740000,Open:4.154,Close:4.153,High:4.154,Low:4.153,Volume:1040,Amount:10400}

	var rs []*entity.Record

	h := NewM1Smoother(security, &r0)
	rs = h.Feed(&r2)
	for _, r := range rs {
		fmt.Printf("%+v\n", r)
	}

	fmt.Println("===================================")

	h1 := NewM1Smoother(security, &r1)
	rs = h1.Feed(&r2)
	for _, r := range rs {
		fmt.Printf("%+v\n", r)
	}
	rs = h1.Feed(&r3)
	for _, r := range rs {
		fmt.Printf("%+v\n", r)
	}
}

func Test_M1SmootherCrossDay(t *testing.T) {
	security := entity.ParseSecurityUnsafe("EOSQFUT.OKEX")

	r0 := entity.Record{Date:1534236840000, Open:4.156, Close:4.159, High: 4.159, Low:4.156, Volume:40, Amount:400}
	r1 := entity.Record{Date:1534237740000, Open:4.156, Close:4.159, High: 4.159, Low:4.156, Volume:40, Amount:400}

	var rs []*entity.Record

	h := NewM1Smoother(security, &r0)
	rs = h.Feed(&r1)
	for _, r := range rs {
		fmt.Printf("%+v\n", r)
	}
}
