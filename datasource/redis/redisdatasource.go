package redisdatasource

import (
	. "github.com/stephenlyu/tds/entity"
	. "github.com/stephenlyu/tds/period"
	"github.com/stephenlyu/tds/datasource"
	"github.com/stephenlyu/tds/storage"
	"fmt"
	"math"
	"sort"
)

type _RedisDataSource struct {
	redisPool *storage.RedisPool
}

func NewRedisDataSource(redisUrl string, redisPassword string) datasource.BaseDataSource {
	redisPool := storage.NewRedisPool(redisUrl, redisPassword)
	return &_RedisDataSource{redisPool: redisPool}
}

func (this *_RedisDataSource) key(security *Security, period Period) string {
	return fmt.Sprintf("%s.%s", period.ShortName(), security.String())
}

func (this *_RedisDataSource) GetData(security *Security, period Period) (error, []Record) {
	key := this.key(security, period)

	array, err := this.redisPool.SortedSetRangeByScore(key, 0, math.MaxUint64, 0, 0)
	if err != nil {
		return err, nil
	}

	ret := make([]Record, len(array))
	for i, str := range array {
		r, err := RecordFromProtoBytes([]byte(str))
		if err != nil {
			return err, nil
		}
		ret[i] = *r
	}

	return nil, ret
}

func (this *_RedisDataSource) GetDataEx(security *Security, period Period, startDate uint64, count int) (error, []Record) {
	key := this.key(security, period)

	array, err := this.redisPool.SortedSetRangeByScore(key, startDate, math.MaxUint64, 0, count)
	if err != nil {
		return err, nil
	}

	ret := make([]Record, len(array))
	for i, str := range array {
		r, err := RecordFromProtoBytes([]byte(str))
		if err != nil {
			return err, nil
		}
		ret[i] = *r
	}

	return nil, ret
}

func (this *_RedisDataSource) GetRangeData(security *Security, period Period, startDate, endDate uint64) (error, []Record) {
	key := this.key(security, period)
	if endDate == 0 {
		endDate = math.MaxUint64
	}

	array, err := this.redisPool.SortedSetRangeByScore(key, startDate, endDate, 0, 0)
	if err != nil {
		return err, nil
	}

	ret := make([]Record, len(array))
	for i, str := range array {
		r, err := RecordFromProtoBytes([]byte(str))
		if err != nil {
			return err, nil
		}
		ret[i] = *r
	}

	return nil, ret
}

func (this *_RedisDataSource) GetDataFromLast(security *Security, period Period, endDate uint64, count int) (error, []Record) {
	if endDate == 0 {
		endDate = math.MaxUint64
	}

	key := this.key(security, period)

	array, err := this.redisPool.SortedSetRevRangeByScore(key, 0, endDate, 0, count)
	if err != nil {
		return err, nil
	}

	ret := make([]Record, len(array))
	for i, str := range array {
		r, err := RecordFromProtoBytes([]byte(str))
		if err != nil {
			return err, nil
		}
		ret[i] = *r
	}

	sort.SliceStable(ret, func (i, j int) bool {
		return ret[i].Date < ret[j].Date
	})

	return nil, ret
}

func (this *_RedisDataSource) GetLastRecord(security *Security, period Period) (error, *Record) {
	key := this.key(security, period)

	array, err := this.redisPool.SortedSetRevRangeByScore(key, 0, math.MaxUint64, 0, 1)
	if err != nil {
		return err, nil
	}

	if len(array) == 0 {
		return nil, nil
	}

	ret, err := RecordFromProtoBytes([]byte(array[0]))
	if err != nil {
		return err, nil
	}

	return nil, ret
}

func (this *_RedisDataSource) AppendData(security *Security, period Period, data []Record) error {
	return this.SaveData(security, period, data)
}

func (this *_RedisDataSource) SaveData(security *Security, period Period, data []Record) error {
	key := this.key(security, period)

	for i := range data {
		r := &data[i]
		bytes, err := r.ToProtoBytes()
		if err != nil {
			return err
		}
		err = this.redisPool.SortedSetRemoveByScore(key, r.Date, r.Date)
		if err != nil {
			return err
		}

		err = this.redisPool.SortedSetAdd(key, bytes, r.Date)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *_RedisDataSource) RemoveData(security *Security, period Period, startDate, endDate uint64) error {
	key := this.key(security, period)

	return this.redisPool.SortedSetRemoveByScore(key, startDate, endDate)
}
