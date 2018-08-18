package storage

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
)

var (
	DEFAULT_REDIS_SERVER = "localhost:6379"
)

type RedisPool struct {
	redis.Pool
}

func NewRedisPool(address string, password string) *RedisPool {
	if address == "" {
		address = DEFAULT_REDIS_SERVER
	}

	return &RedisPool{
		Pool: redis.Pool {
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", address)
				if err != nil {
					return nil, err
				}
				if password != "" {
					if _, err := c.Do("AUTH", password); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
}


// SortedSet

func (redisPool *RedisPool) SortedSetAdd(key string, data []byte, score uint64) error {
	c := redisPool.Get()
	defer c.Close()

	_, err := c.Do("ZADD", key, score, data)
	return err
}

// min, max: inclusive
func (redisPool *RedisPool) SortedSetRangeByScore(key string, min, max uint64, offset, count int) ([]string, error) {
	c := redisPool.Get()
	defer c.Close()

	args := []interface{}{key, min, max}

	if count > 0 {
		args = append(args, "LIMIT")
		args = append(args, offset)
		args = append(args, count)
	}

	fmt.Printf("%+v\n", args)
	values, err := redis.Strings(c.Do("ZRANGEBYSCORE", args...))

	if err != nil && err.Error() == "redigo: nil returned" {
		return []string{}, nil
	}
	return values, err
}

// min, max: inclusive
func (redisPool *RedisPool) SortedSetRevRangeByScore(key string, min, max uint64, offset, count int) ([]string, error) {
	c := redisPool.Get()
	defer c.Close()

	args := []interface{}{key, max, min}

	if count > 0 {
		args = append(args, "LIMIT")
		args = append(args, offset)
		args = append(args, count)
	}

	values, err := redis.Strings(c.Do("ZREVRANGEBYSCORE", args...))

	if err != nil && err.Error() == "redigo: nil returned" {
		return []string{}, nil
	}
	return values, err
}

// min, max: inclusive
func (redisPool *RedisPool) SortedSetRemoveByScore(key string, min, max uint64) error {
	c := redisPool.Get()
	defer c.Close()

	_, err := c.Do("ZREMRANGEBYSCORE", key, min, max)
	return err
}
