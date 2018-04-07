package periodmanager

import (
	. "github.com/stephenlyu/tds/period"
	"sync"
	"github.com/stephenlyu/tds/util"
	"sort"
)

type PeriodManager interface {
	// 检查PeriodManager中是否包含period
	// @param: period - A period
	HasPeriod(period Period) bool

	// 添加周期到PeriodManager
	AddPeriod(period Period)

	// 检查是否是basic period
	IsBasicPeriod(period Period) bool

	// 获取period的依赖周期列表
	GetPeriodDependencies(period Period) []Period

	// 获取所有添加周期的按照Merge次序排序的周期列表
	GetOrderedPeriods() []Period
}

type defaultPeriodManager struct {
	basicPeriods []Period
	periods 	[]Period

	orderedPeriods []Period

	periodDependencies map[string][]Period

	lock sync.RWMutex
}

func NewDefaultPeriodManager(basicPeriods []Period) PeriodManager {
	basicPeriods = append([]Period{}, basicPeriods...)
	sort.SliceStable(basicPeriods, func (i, j int) bool {
		return basicPeriods[i].Lt(basicPeriods[j])
	})
	return &defaultPeriodManager{
		basicPeriods: basicPeriods,
		periodDependencies: make(map[string][]Period),
		orderedPeriods: append([]Period{}, basicPeriods...),
	}
}


func (this *defaultPeriodManager) HasPeriod(period Period) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, p := range this.basicPeriods {
		if p.ShortName() == period.ShortName() {
			return true
		}
	}
	for _, p := range this.periods {
		if p.ShortName() == period.ShortName() {
			return true
		}
	}

	return false
}

func (this *defaultPeriodManager) AddPeriod(period Period) {
	if this.HasPeriod(period) {
		return
	}

	this.lock.RLock()
	dependencies, ok := this.periodDependencies[period.ShortName()]
	this.lock.RUnlock()
	if !ok {
		dependencies = this.GetPeriodDependencies(period)
	}

	dependencies = append(dependencies, period)
	this.mergeOrderedPeriods(dependencies)

	this.lock.Lock()
	this.periods = append(this.periods, period)
	this.lock.Unlock()
}

func (this *defaultPeriodManager) IsBasicPeriod(period Period) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, p := range this.basicPeriods {
		if p.Eq(period) {
			return true
		}
	}
	return false
}

func (this *defaultPeriodManager) getPeriodDependenciesInternal(period Period) []Period {
	var ret Period
	// Find max period from custom periods
	for _, p := range this.periods {
		if period.CanConvertFrom(p) {
			if ret == nil || ret.Lt(p) {
				ret = p
			}
		}
	}
	if ret != nil {
		if ret.Eq(period) {
			return nil
		}
		return []Period{ret}
	}

	// Find max period from basic periods
	for _, p := range this.basicPeriods {
		if period.CanConvertFrom(p) {
			if ret == nil || ret.Lt(p) {
				ret = p
			}
		}
	}
	if ret != nil {
		if ret.Eq(period) {
			return nil
		}
		return []Period{ret}
	}

	periods := append(this.getPeriodDependenciesInternal(period.BasicMergePeriod()), period.BasicMergePeriod())

	return periods
}

func (this *defaultPeriodManager) GetPeriodDependencies(period Period) []Period {
	this.lock.RLock()
	periods := this.getPeriodDependenciesInternal(period)
	this.lock.RUnlock()

	this.lock.Lock()
	this.periodDependencies[period.ShortName()] = periods
	this.lock.Unlock()

	return periods
}

func (this *defaultPeriodManager) mergeOrderedPeriods(periods []Period) {
	this.lock.Lock()
	defer this.lock.Unlock()

	var orderPeriods []Period

	i, j := 0, 0
	for ; i < len(this.orderedPeriods) && j < len(periods); {
		p1 := this.orderedPeriods[i]
		p2 := periods[j]
		switch {
		case p1.Eq(p2):
			orderPeriods = append(orderPeriods, p1)
			i++
			j++
		case p1.Lt(p2):
			orderPeriods = append(orderPeriods, p1)
			i++
		default:
			orderPeriods = append(orderPeriods, p2)
			j++
		}
	}

	if i < len(this.orderedPeriods) {
		util.Assert(j >= len(periods), "")
		orderPeriods = append(orderPeriods, this.orderedPeriods[i:]...)
	}

	if j < len(periods) {
		util.Assert(i >= len(this.orderedPeriods), "")
		orderPeriods = append(orderPeriods, periods[j:]...)
	}

	this.orderedPeriods = orderPeriods
}

func (this *defaultPeriodManager) GetOrderedPeriods() []Period {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return append([]Period{}, this.orderedPeriods...)
}
