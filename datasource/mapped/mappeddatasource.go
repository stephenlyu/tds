package mappeddatasource

import (
	. "github.com/stephenlyu/tds/datasource"
	. "github.com/stephenlyu/tds/entity"
	. "github.com/stephenlyu/tds/period"
	"sort"
	"math"
	"github.com/stephenlyu/tds/util"
)

type _MappedDataSource struct {
	mapper DateRangeMapper
	targetDs BaseDataSource
}

func NewMapperDataSource() MappedDataSource {
	return &_MappedDataSource{}
}

func (this *_MappedDataSource) SetMapper(mapper DateRangeMapper) {
	this.mapper = mapper
}

func (this *_MappedDataSource) SetTargetDataSource(ds BaseDataSource) {
	this.targetDs = ds
}

func (this *_MappedDataSource) getRanges(security *Security) []DateRange {
	ranges := this.mapper.MapDateRanges(security)
	sort.SliceStable(ranges, func (i, j int) bool {
		return ranges[i].StartDate < ranges[j].StartDate
	})
	return ranges
}

func (this *_MappedDataSource) GetData(security *Security, period Period) (error, []Record) {
	if this.mapper == nil || this.targetDs == nil {
		panic("IllegalStateException")
	}
	ranges := this.getRanges(security)
	var ret []Record

	for _, r := range ranges {
		err, data := this.targetDs.GetRangeData(r.Security, period, r.StartDate, r.EndDate - 1)
		if err != nil {
			return err, nil
		}
		ret = append(ret, data...)
	}
	return nil, ret
}

func (this *_MappedDataSource) GetDataEx(security *Security, period Period, startDate uint64, count int) (error, []Record) {
	panic("Unimplemented")
	return
}

func (this *_MappedDataSource) GetRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record) {
	if this.mapper == nil || this.targetDs == nil {
		panic("IllegalStateException")
	}

	if endDate == 0 {
		endDate = uint64(math.MaxInt64)
	}

	ranges := this.getRanges(security)
	var ret []Record

	for _, r := range ranges {
		if r.EndDate == 0 {
			r.EndDate = uint64(math.MaxInt64)
		}
		start := util.MaxUInt64(r.StartDate, startDate)
		end := util.MinUInt64(r.EndDate-1, endDate)

		err, data := this.targetDs.GetRangeData(r.Security, period, start, end)
		if err != nil {
			return err, nil
		}
		ret = append(ret, data...)
	}
	return nil, ret
}

func (this *_MappedDataSource) GetDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record) {
	panic("Unimplemented")
	return
}

func (this *_MappedDataSource) GetLastRecord(security *Security, period Period) (error, *Record) {
	panic("Unimplemented")
	return
}

func (this *_MappedDataSource) AppendData(security *Security, period Period, data []Record) error {
	panic("Unimplemented")
	return nil
}

func (this *_MappedDataSource) SaveData(security *Security, period Period, data []Record) error {
	panic("Unimplemented")
	return nil
}

func (this *_MappedDataSource) RemoveData(security *Security, period Period, startDate, endDate uint64) error {
	panic("Unimplemented")
	return nil
}
