package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// returns a new redis cache.
func NewRedisClient(address, passwordKey string) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: os.Getenv(passwordKey),
		DB:       0,
		// TODO use config
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	})

	ctx := context.Background()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("rdb.Ping: %w", err)
	}

	return redisClient, nil
}
