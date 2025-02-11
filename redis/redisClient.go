package redis

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	RedisInstance *redis.Client
	once          sync.Once
)

func NewRedisInstance() *redis.Client {
	once.Do(func() {
		RedisInstance = redis.NewClient(&redis.Options{
			Addr:         "localhost:6379",
			Password:     "",
			DB:           0,
			PoolSize:     5,
			MinIdleConns: 2,
			PoolTimeout:  30,
		})
	})
	return RedisInstance
}
