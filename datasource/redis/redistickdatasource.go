package redisdatasource

import (
	. "github.com/stephenlyu/tds/entity"
	"github.com/stephenlyu/tds/datasource"
	"github.com/stephenlyu/tds/storage"
	"fmt"
	"math"
)

type _RedisTickDataSource struct {
	redisPool *storage.RedisPool
}

func NewTickRedisDataSource(redisUrl string, redisPassword string) datasource.TickDataSource {
	redisPool := storage.NewRedisPool(redisUrl, redisPassword)
	return &_RedisTickDataSource{redisPool: redisPool}
}

func (this *_RedisTickDataSource) key(security *Security) string {
	return fmt.Sprintf("tick.%s", security.String())
}

func (this *_RedisTickDataSource) Get(security *Security, startTs, endTs uint64) (error, []TickItem) {
	key := this.key(security)

	array, err := this.redisPool.SortedSetRangeByScore(key, 0, math.MaxUint64, 0, 0)
	if err != nil {
		return err, nil
	}

	ret := make([]TickItem, len(array))
	for i, str := range array {
		r, err := TickItemFromProtoBytes([]byte(str))
		if err != nil {
			return err, nil
		}
		ret[i] = *r
	}

	return nil, ret
}

func (this *_RedisTickDataSource) Remove(security *Security, startTs, endTs uint64) error {
	key := this.key(security)

	return this.redisPool.SortedSetRemoveByScore(key, startTs, endTs)
}

func (this *_RedisTickDataSource) Save(security *Security, ticks []TickItem) error {
	key := this.key(security)

	for i := range ticks {
		r := &ticks[i]
		bytes, err := r.ToProtoBytes()
		if err != nil {
			return err
		}
		err = this.redisPool.SortedSetAdd(key, bytes, r.Timestamp)
		if err != nil {
			return err
		}
	}
	return nil
}
