package shared

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis(addr string) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
	_, err := RedisClient.Ping(ctx).Result()
	return err
}

func SetCache(key string, value string, ttl time.Duration) error {
	return RedisClient.Set(ctx, key, value, ttl).Err()
}

func GetCache(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}
