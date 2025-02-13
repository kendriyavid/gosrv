package redis

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	RedisInstance *redis.Client
	once          sync.Once
)

func NewRedisInstance() *redis.Client {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system environment variables")
		}

		RedisInstance = redis.NewClient(&redis.Options{
			Addr:         os.Getenv("REDIS_HOST"),
			Password:     os.Getenv("REDIS_PASSWORD"),
			DB:           0,
			PoolSize:     5,
			MinIdleConns: 2,
			PoolTimeout:  30 * time.Second,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := RedisInstance.Ping(ctx).Err(); err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}
	})
	return RedisInstance
}
