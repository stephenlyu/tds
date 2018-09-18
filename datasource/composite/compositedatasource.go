package compositedatasource

import (
	. "github.com/stephenlyu/tds/datasource"
	. "github.com/stephenlyu/tds/entity"
	. "github.com/stephenlyu/tds/period"
	"fmt"
	"github.com/stephenlyu/tds/date"
)

type _CompositeDataSource struct {
	subDataSources []BaseDataSource
}

func NewCompositeDataSource() CompositeDataSource {
	return &_CompositeDataSource{}
}

func (this *_CompositeDataSource) AddSubDatasource(ds BaseDataSource) {
	this.subDataSources = append(this.subDataSources, ds)
}

func (this *_CompositeDataSource) GetData(security *Security, period Period) (error, []Record) {
	var ret []Record
	var startDate uint64
	for _, ds := range this.subDataSources {
		err, data := ds.GetRangeData(security, period, startDate, 0)
		if err != nil {
			return err, nil
		}
		ret = append(ret, data...)
		if len(ret) > 0 {
			startDate = ret[len(ret) - 1].Date + 1 // 数据周期最小为分钟，分钟毫秒数+1不会导致结果少数据
		}
	}
	return nil, ret
}

func (this *_CompositeDataSource) GetDataEx(security *Security, period Period, startDate uint64, count int) (error, []Record) {
	var ret []Record
	n := count
	for _, ds := range this.subDataSources {
		err, data := ds.GetDataEx(security, period, startDate, n)
		if err != nil {
			return err, nil
		}
		ret = append(ret, data...)
		if len(ret) > 0 {
			startDate = ret[len(ret) - 1].Date + 1
		}
		n -= len(data)
		if n <= 0 {
			break
		}
	}
	if len(ret) > count {
		ret = ret[:count]
	}
	return nil, ret
}

func (this *_CompositeDataSource) GetRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record) {
	var ret []Record
	for _, ds := range this.subDataSources {
		err, data := ds.GetRangeData(security, period, startDate, endDate)
		if err != nil {
			return err, nil
		}
		fmt.Println(date.Timestamp2SecondString(startDate), date.Timestamp2SecondString(endDate), len(data))

		ret = append(ret, data...)
		if len(ret) > 0 {
			startDate = ret[len(ret) - 1].Date + 1
		}
	}
	return nil, ret
}

func (this *_CompositeDataSource) GetDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record) {
	var ret []Record
	n := count
	for i := len(this.subDataSources) - 1; i >= 0; i-- {
		ds := this.subDataSources[i]
		err, data := ds.GetDataFromLast(security, period, endDate, n)
		if err != nil {
			return err, nil
		}
		ret = append(ret, data...)
		if len(ret) > 0 {
			endDate = ret[0].Date - 1 // endData is inclusive
		}
		n -= len(data)
		if n <= 0 {
			break
		}
	}
	if len(ret) > count {
		ret = ret[:count]
	}
	return nil, ret
}

func (this *_CompositeDataSource) GetLastRecord(security *Security, period Period) (error, *Record) {
	for i := len(this.subDataSources) - 1; i >= 0; i-- {
		ds := this.subDataSources[i]
		err, r := ds.GetLastRecord(security, period)
		if err != nil {
			return err, nil
		}
		if r != nil {
			return nil, r
		}
	}
	return nil, nil
}

func (this *_CompositeDataSource) AppendData(security *Security, period Period, data []Record) error {
	panic("Unimplemented")
	return nil
}

func (this *_CompositeDataSource) SaveData(security *Security, period Period, data []Record) error {
	panic("Unimplemented")
	return nil
}

func (this *_CompositeDataSource) RemoveData(security *Security, period Period, startDate, endDate uint64) error {
	panic("Unimplemented")
	return nil
}
