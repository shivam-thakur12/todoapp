package todo

import (
	"context"
	"testing"
	"time"

	rediss "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache_Set(t *testing.T) {
	// Initialize the Redis client
	client := rediss.NewClient(&rediss.Options{
		Addr: "localhost:6379", // Adjust according to your Redis setup
	})

	// Create a RedisCache instance
	cache := NewRedisCache(client)

	key := "test-key"
	value := "test-value"
	expiration := time.Minute

	// Test the Set function
	err := cache.Set(key, value, expiration)
	assert.NoError(t, err)

	// Verify by checking if the value was set correctly
	gotValue, err := client.Get(context.Background(), key).Result()
	assert.NoError(t, err)
	assert.Equal(t, value, gotValue)
}

func TestRedisCache_Get(t *testing.T) {
	// Initialize the Redis client
	client := rediss.NewClient(&rediss.Options{
		Addr: "localhost:6379", // Adjust according to your Redis setup
	})

	// Create a RedisCache instance
	cache := NewRedisCache(client)

	key := "test-key"
	expectedValue := "test-value"

	// Set the value first to ensure it exists
	err := cache.Set(key, expectedValue, time.Minute)
	assert.NoError(t, err)

	// Test the Get function
	gotValue, err := cache.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, gotValue)
}

func TestRedisCache_Del(t *testing.T) {
	// Initialize the Redis client
	client := rediss.NewClient(&rediss.Options{
		Addr: "localhost:6379", // Adjust according to your Redis setup
	})

	// Create a RedisCache instance
	cache := NewRedisCache(client)

	key := "test-key"
	value := "test-value"

	// Set the value first to ensure it exists
	err := cache.Set(key, value, time.Minute)
	assert.NoError(t, err)

	// Test the Del function
	err = cache.Del(key)
	assert.NoError(t, err)

	// Verify by checking if the key was deleted
	_, err = client.Get(context.Background(), key).Result()
	assert.Error(t, err)
	assert.Equal(t, rediss.Nil, err) // redis.Nil means the key does not exist
}
