package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	client *redis.Client
}

//returns a new redis cache
func NewRedisCache(address, password string) (*redisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("rdb.Ping: %w", err)
	}

	return &redisCache{
		client: rdb,
	}, nil
}

func (cache *redisCache) Ping() error {
	ctx := context.Background()
	_, err := cache.client.Ping(ctx).Result()
	return err
}

func (cache *redisCache) ReadString(key string) (string, error) {
	val, err := cache.client.Get(context.Background(), key).Result()
	return val, err
}

func (cache *redisCache) WriteString(key, value string) error {
	_, err := cache.client.Set(context.Background(), key, value, 0).Result()
	return err
}

func (cache *redisCache) WriteWithExpiry(key string, value string, expiryTime time.Duration) error {
	_, err := cache.client.SetEX(context.Background(), key, value, expiryTime).Result()
	return err
}

func (cache *redisCache) Increment(key string) (int64, error) {
	val, err := cache.client.Incr(context.Background(), key).Result()
	return val, err
}

func (cache *redisCache) AddExpiry(key string, expiryTime time.Duration) error {
	_, err := cache.client.Expire(context.Background(), key, expiryTime).Result()
	return err
}

func (cache *redisCache) Delete(key string) error {
	_, err := cache.client.Del(context.Background(), key).Result()
	return err
}

func (cache *redisCache) Close() {
	cache.client.Close()
}
