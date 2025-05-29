package cache

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Cache struct {
	Redis *redis.Client
}

func NewCache(r *redis.Client) *Cache {
	return &Cache{Redis: r}
}

func (c *Cache) GetSearchResultsHashKey(city string, min, max float64, category string) string {
	keyData := fmt.Sprintf("%s:%f:%f:%s", city, min, max, category)
	hash := sha1.Sum([]byte(keyData))
	return fmt.Sprintf("search:%x", hash)
}

func (c *Cache) Save(key string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.Redis.Set(ctx, key, jsonData, 10*time.Minute).Err()
}

func (c *Cache) Load(key string, out interface{}) error {
	val, err := c.Redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), out)
}
