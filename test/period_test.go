package test


import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stephenlyu/tds"
)

var _ = Describe("TestPeriod", func() {
	It("test", func (){
		var p, p1 tds.Period
		var err error

		err, p = tds.PeriodFromString("M1")
		Expect(err).To(BeNil())
		Expect(p).To(Not(BeNil()))
		Expect(p.ShortName()).To(Equal("M1"))
		Expect(p.Name()).To(Equal("MINUTE1"))

		err, p = tds.PeriodFromString("M")
		Expect(err).To(Not(BeNil()))

		err, p = tds.PeriodFromString("M1")
		Expect(err).To(BeNil())

		err, p1 = tds.PeriodFromString("M5")
		Expect(err).To(BeNil())

		Expect(p.Lt(p1)).To(Equal(true))
		Expect(p1.Gt(p)).To(Equal(true))

		err, p1 = tds.PeriodFromString("M1")
		Expect(err).To(BeNil())
		Expect(p.Eq(p1)).To(Equal(true))

		err, p1 = tds.PeriodFromString("D1")
		Expect(err).To(BeNil())

		Expect(p.Lt(p1)).To(Equal(true))
		Expect(p1.Gt(p)).To(Equal(true))


		var pm, pd, pw, pn, pq, py tds.Period
		var pm1, pd1, pw1, pn1, pq1, py1 tds.Period

		err, pm = tds.PeriodFromString("M1")
		Expect(err).To(BeNil())
		err, pm1 = tds.PeriodFromString("M5")
		Expect(err).To(BeNil())

		err, pd = tds.PeriodFromString("D1")
		Expect(err).To(BeNil())
		err, pd1 = tds.PeriodFromString("D5")
		Expect(err).To(BeNil())

		err, pw = tds.PeriodFromString("W1")
		Expect(err).To(BeNil())
		err, pw1 = tds.PeriodFromString("W5")
		Expect(err).To(BeNil())

		err, pn = tds.PeriodFromString("N1")
		Expect(err).To(BeNil())
		err, pn1 = tds.PeriodFromString("N5")
		Expect(err).To(BeNil())

		err, pq = tds.PeriodFromString("Q1")
		Expect(err).To(BeNil())
		err, pq1 = tds.PeriodFromString("Q5")
		Expect(err).To(BeNil())

		err, py = tds.PeriodFromString("Y1")
		Expect(err).To(BeNil())
		err, py1 = tds.PeriodFromString("Y5")
		Expect(err).To(BeNil())

		Expect(pm.CanConvertTo(pm1)).To(Equal(true))
		Expect(pm.CanConvertTo(pd)).To(Equal(true))
		Expect(pm.CanConvertTo(pd1)).To(Equal(false))
		Expect(pm.CanConvertTo(pw)).To(Equal(false))
		Expect(pm.CanConvertTo(pn)).To(Equal(false))
		Expect(pm.CanConvertTo(pq)).To(Equal(false))
		Expect(pm.CanConvertTo(py)).To(Equal(false))

		Expect(pd.CanConvertTo(pm)).To(Equal(false))
		Expect(pd.CanConvertTo(pd1)).To(Equal(true))
		Expect(pd.CanConvertTo(pw)).To(Equal(true))
		Expect(pd.CanConvertTo(pn)).To(Equal(true))
		Expect(pd.CanConvertTo(pq)).To(Equal(true))
		Expect(pd.CanConvertTo(py)).To(Equal(true))

		Expect(pw.CanConvertTo(pm)).To(Equal(false))
		Expect(pw.CanConvertTo(pd)).To(Equal(false))
		Expect(pw.CanConvertTo(pw1)).To(Equal(true))
		Expect(pw.CanConvertTo(pn)).To(Equal(false))
		Expect(pw.CanConvertTo(pq)).To(Equal(false))
		Expect(pw.CanConvertTo(py)).To(Equal(false))

		Expect(pn.CanConvertTo(pm)).To(Equal(false))
		Expect(pn.CanConvertTo(pd)).To(Equal(false))
		Expect(pn.CanConvertTo(pw)).To(Equal(false))
		Expect(pn.CanConvertTo(pn1)).To(Equal(true))
		Expect(pn.CanConvertTo(pq)).To(Equal(true))
		Expect(pn.CanConvertTo(py)).To(Equal(true))

		Expect(pq.CanConvertTo(pm)).To(Equal(false))
		Expect(pq.CanConvertTo(pd)).To(Equal(false))
		Expect(pq.CanConvertTo(pw)).To(Equal(false))
		Expect(pq.CanConvertTo(pn)).To(Equal(false))
		Expect(pq.CanConvertTo(pq1)).To(Equal(true))
		Expect(pq.CanConvertTo(py)).To(Equal(true))

		Expect(py.CanConvertTo(pm)).To(Equal(false))
		Expect(py.CanConvertTo(pd)).To(Equal(false))
		Expect(py.CanConvertTo(pw)).To(Equal(false))
		Expect(py.CanConvertTo(pn)).To(Equal(false))
		Expect(py.CanConvertTo(pq)).To(Equal(false))
		Expect(py.CanConvertTo(py1)).To(Equal(true))

		Expect(py.CanConvertTo(py1)).To(Equal(py1.CanConvertFrom(py)))
	})
})
