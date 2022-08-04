package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// returns a new redis cache.
func NewRedisClient(address, password string) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("rdb.Ping: %w", err)
	}

	return redisClient, nil
}
