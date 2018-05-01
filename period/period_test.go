package period_test

import (
	"testing"
	"github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/util"
	"fmt"
	"unsafe"
)

func assert(cond bool, msg string) {
	util.Assert(cond, msg)
}

type typeAlg struct {
	// function for hashing objects of this type
	// (ptr to object, seed) -> hash
	hash func(unsafe.Pointer, uintptr) uintptr
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
}

type tflag uint8
type nameOff int32
type typeOff int32

type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldalign uint8
	kind       uint8
	alg        *typeAlg
					   // gcdata stores the GC type data for the garbage collector.
					   // If the KindGCProg bit is set in kind, gcdata is a GC program.
					   // Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}

type imethod struct {
	name nameOff
	ityp typeOff
}

type name struct {
	bytes *byte
}

type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type eface struct {
	_type *_type
	data  unsafe.Pointer
}

type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod
}

type itab struct {
	inter  *interfacetype
	_type  *_type
	link   *itab
	hash   uint32 // copy of _type.hash. Used for type switches.
	bad    bool   // type does not implement interface
	inhash bool   // has this itab been added to hash?
	unused [2]byte
	fun    [1]uintptr // variable sized
}

type hmap struct {
							  // Note: the format of the Hmap is encoded in ../../cmd/internal/gc/reflect.go and
							  // ../reflect/type.go. Don't change this structure without also changing that code!
	count      int            // # live cells == size of map.  Must be first (used by len() builtin)
	flags      uint8
	B          uint8          // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow  uint16         // approximate number of overflow buckets; see incrnoverflow for details
	hash0      uint32         // hash seed

	buckets    unsafe.Pointer // array of 2^B Buckets. may be nil if count==0.
	oldbuckets unsafe.Pointer // previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)
}

func TestPeriod(t *testing.T) {
	var p, p1 period.Period
	var err error
	//
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

	assert(!period.PeriodFromStringUnsafe("Q1").CanConvertFrom(period.PeriodFromStringUnsafe("D3")), "")

	var m = make(map[period.Period]period.Period)

	err, p2 := period.PeriodFromString("M1")

	mm := (*hmap)(unsafe.Pointer(&m))

	i1 := (*iface)(unsafe.Pointer(&p))
	i2 := (*iface)(unsafe.Pointer(&p2))

	fmt.Println("i1.hash:", i1.tab._type.alg.hash(i1.data, uintptr(mm.hash0)))
	fmt.Println("i2.hash:", i2.tab._type.alg.hash(i2.data, uintptr(mm.hash0)))

	m[p] = p1
	fmt.Println(p, m[p], p1)
	fmt.Println(p, m[p2], p1)
	//assert(m[p] == p1, "")
}
