package test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stephenlyu/tds"
)

var _ = Describe("TestRecord", func() {
	It("test", func (){
		r := &tds.Record{
			Date: uint32(20170704),
			Open: 8500,
			Close: 9050,
			High: 9190,
			Low: 8460,
			Volume: 10000,
			Amount: 1000000,
		}

		r1 := &tds.Record{}
		err := tds.RecordFromBytes(r.Bytes(), r1)
		Expect(err).To(BeNil())
		Expect(r1.Date).To(Equal(r.Date))
		Expect(r1.Open).To(Equal(r.Open))
		Expect(r1.Close).To(Equal(r.Close))
		Expect(r1.High).To(Equal(r.High))
		Expect(r1.Low).To(Equal(r.Low))
		Expect(r1.Volume).To(Equal(r.Volume))
		Expect(r1.Amount).To(Equal(r.Amount))
	})
})
