package shared

import (
	"github.com/redis/go-redis/v9"
)

func InitRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
}
