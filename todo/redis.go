package todo

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache defines the interface for Redis operations
type RedisCache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
}

// redisCache is the implementation of the RedisCache interface
type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new RedisCache instance
func NewRedisCache(client *redis.Client) RedisCache {
	return &redisCache{
		client: client,
	}
}

var ctx = context.Background()

// Set adds a key-value pair to the cache
func (r *redisCache) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key from the cache
func (r *redisCache) Get(key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del removes a key from the cache
func (r *redisCache) Del(key string) error {
	return r.client.Del(ctx, key).Err()
}
