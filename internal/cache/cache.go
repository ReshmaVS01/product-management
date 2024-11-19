package cache

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func Set(key, value string) {
	err := RedisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}
}

func Get(key string) (string, error) {
	val, err := RedisClient.Get(ctx, key).Result()
	return val, err
}

