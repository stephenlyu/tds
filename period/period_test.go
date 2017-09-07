package period_test

import (
	"testing"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/util"
)

func assert(cond bool, msg string) {
	util.Assert(cond, msg)
}

func TestPeriod(t *testing.T) {
	var p, p1 period.Period
	var err error

	err, p = period.PeriodFromString("M1")
	assert(err == nil, "err == nil")
	assert(p != nil, "p != nil")
	assert(p.ShortName() == "M1", `p.ShortName() == "M1"`)
	assert(p.Name() == "MINUTE1", `p.ShortName() == "MINUTE1"`)

	err, p = period.PeriodFromString("M")
	assert(err != nil, "err != nil")

	err, p = period.PeriodFromString("M1")
	assert(err == nil, "err == nil")

	err, p1 = period.PeriodFromString("M5")
	assert(err == nil, "err == nil")

	assert(p.Lt(p1), "p.Lt(p1)")
	assert(p1.Gt(p), "p1.Gt(p)")

	err, p1 = period.PeriodFromString("M1")
	assert(err == nil, "err == nil")
	assert(p.Eq(p1), "p.Eq(p1)")

	err, p1 = period.PeriodFromString("D1")
	assert(err == nil, "err == nil")

	assert(p.Lt(p1), "p.Lt(p1)")
	assert(p1.Gt(p), "p1.Gt(p)")

	var pm, pd, pw, pn, pq, py period.Period
	var pm1, pd1, pw1, pn1, pq1, py1 period.Period

	err, pm = period.PeriodFromString("M1")
	assert(err == nil, "err == nil")
	err, pm1 = period.PeriodFromString("M5")
	assert(err == nil, "err == nil")

	err, pd = period.PeriodFromString("D1")
	assert(err == nil, "err == nil")
	err, pd1 = period.PeriodFromString("D5")
	assert(err == nil, "err == nil")

	err, pw = period.PeriodFromString("W1")
	assert(err == nil, "err == nil")
	err, pw1 = period.PeriodFromString("W5")
	assert(err == nil, "err == nil")

	err, pn = period.PeriodFromString("N1")
	assert(err == nil, "err == nil")
	err, pn1 = period.PeriodFromString("N5")
	assert(err == nil, "err == nil")

	err, pq = period.PeriodFromString("Q1")
	assert(err == nil, "err == nil")
	err, pq1 = period.PeriodFromString("Q5")
	assert(err == nil, "err == nil")

	err, py = period.PeriodFromString("Y1")
	assert(err == nil, "err == nil")
	err, py1 = period.PeriodFromString("Y5")
	assert(err == nil, "err == nil")

	assert(pm.CanConvertTo(pm1), "pm.CanConvertTo(pm1)")
	assert(pm.CanConvertTo(pd), "pm.CanConvertTo(pd)")
	assert(!pm.CanConvertTo(pd1), "!pm.CanConvertTo(pd1)")
	assert(!pm.CanConvertTo(pw), "!pm.CanConvertTo(pw)")
	assert(!pm.CanConvertTo(pn), "!pm.CanConvertTo(pn)")
	assert(!pm.CanConvertTo(pq), "!pm.CanConvertTo(pq)")
	assert(!pm.CanConvertTo(py), "!pm.CanConvertTo(py)")

	assert(!pd.CanConvertTo(pm), "!pd.CanConvertTo(pm)")
	assert(pd.CanConvertTo(pd1), "pd.CanConvertTo(pd1)")
	assert(pd.CanConvertTo(pw), "pd.CanConvertTo(pw)")
	assert(pd.CanConvertTo(pn), "pd.CanConvertTo(pn)")
	assert(pd.CanConvertTo(pq), "pd.CanConvertTo(pq)")
	assert(pd.CanConvertTo(py), "pd.CanConvertTo(py)")

	assert(!pw.CanConvertTo(pm), "!pw.CanConvertTo(pm)")
	assert(!pw.CanConvertTo(pd), "!pw.CanConvertTo(pd)")
	assert(pw.CanConvertTo(pw1), "pw.CanConvertTo(pw1)")
	assert(!pw.CanConvertTo(pn), "!pw.CanConvertTo(pn)")
	assert(!pw.CanConvertTo(pq), "!pw.CanConvertTo(pq)")
	assert(!pw.CanConvertTo(py), "!pw.CanConvertTo(py)")

	assert(!pn.CanConvertTo(pm), "!pn.CanConvertTo(pm)")
	assert(!pn.CanConvertTo(pd), "!pn.CanConvertTo(pd)")
	assert(!pn.CanConvertTo(pw), "!pn.CanConvertTo(pw)")
	assert(pn.CanConvertTo(pn1), "pn.CanConvertTo(pn1)")
	assert(pn.CanConvertTo(pq), "pn.CanConvertTo(pq)")
	assert(pn.CanConvertTo(py), "pn.CanConvertTo(py)")

	assert(!pq.CanConvertTo(pm), "!pq.CanConvertTo(pm)")
	assert(!pq.CanConvertTo(pd), "!pq.CanConvertTo(pd)")
	assert(!pq.CanConvertTo(pw), "!pq.CanConvertTo(pw)")
	assert(!pq.CanConvertTo(pn), "!pq.CanConvertTo(pn)")
	assert(pq.CanConvertTo(pq1), "pq.CanConvertTo(pq1)")
	assert(pq.CanConvertTo(py), "pq.CanConvertTo(py)")

	assert(!py.CanConvertTo(pm), "!py.CanConvertTo(pm)")
	assert(!py.CanConvertTo(pd), "!py.CanConvertTo(pd)")
	assert(!py.CanConvertTo(pw), "!py.CanConvertTo(pw)")
	assert(!py.CanConvertTo(pn), "!py.CanConvertTo(pn)")
	assert(!py.CanConvertTo(pq), "!py.CanConvertTo(pq)")
	assert(py.CanConvertTo(py1), "py.CanConvertTo(py1)")

	assert(py.CanConvertTo(py1) == py1.CanConvertFrom(py), "py.CanConvertTo(py1) == py1.CanConvertFrom(py)")
}