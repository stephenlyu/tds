package test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stephenlyu/tds"
	"fmt"
)

var _ = Describe("TestDate", func() {
	It("test", func (){
		var d tds.Date
		var err error

		err, d = tds.NewDateFromDayString("20170704")
		Expect(err).To(BeNil())
		Expect(d.DayString()).To(Equal("20170704"))
		Expect(d.DayDate()).To(Equal(uint32(20170704)))

		d = tds.NewDateFromDayDate(uint32(20170704))
		Expect(d.DayString()).To(Equal("20170704"))

		err, d = tds.NewDateFromMinuteString("20170704 09:30:00")
		Expect(err).To(BeNil())
		Expect(d.MinuteString()).To(Equal("20170704 09:30:00"))
		fmt.Println(d.MinuteDate(), d.MinuteString())

		d = tds.NewDateFromMinuteDate(d.MinuteDate())
		Expect(d.MinuteString()).To(Equal("20170704 09:30:00"))
	})
})
