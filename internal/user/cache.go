package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kaustubhbabar5/gh-api-client/pkg/github"
	"go.uber.org/zap"
)

const (
	userCachePrefix = "USER_INFO_"
)

type Cache interface {
	// takes multiple users and cache data for give duration.
	// returns back slice of usernames for which caching operation failed and error.
	CacheUserInfo(users github.Users, duration time.Duration) []error
	//
	GetUserInfo(usernames []string) (
		users github.Users,
		notFound []string,
		errs []error,
	)
}

type c struct {
	logger *zap.Logger
	rc     *redis.Client
}

func NewCache(logger *zap.Logger, cache *redis.Client) *c {
	return &c{logger, cache}
}

func (c *c) CacheUserInfo(users github.Users, duration time.Duration) []error {
	ctx := context.Background()
	errs := make([]error, 0)

	for _, user := range users {
		userInfo, err := json.Marshal(user)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to marshal data for %s with error :%w", user.Login, err))
			continue
		}
		err = c.rc.SetEX(ctx, userCachePrefix+user.Login, userInfo, duration).Err()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to cache data for %s with error :%w", user.Login, err))
			continue
		}
	}

	return errs
}

func (c *c) GetUserInfo(usernames []string) ( //nolint: nonamedreturns // TODO add directive
	users github.Users,
	notFound []string,
	errs []error,
) {
	ctx := context.Background()

	redisKeys := make([]string, len(usernames))

	notFound = make([]string, 0)
	users = make(github.Users, 0)
	errs = make([]error, 0)

	for i := range usernames {
		redisKeys[i] = userCachePrefix + usernames[i]
	}

	vals, err := c.rc.MGet(ctx, redisKeys...).Result()
	if err != nil {
		notFound = usernames
		errs = []error{fmt.Errorf("c.rc.MGet :%w", err)}
		return
	}
	for index, val := range vals {
		if val == nil {
			notFound = append(notFound, usernames[index])
			continue
		}

		switch data := val.(type) {
		case string:

			var user github.User
			e := json.Unmarshal([]byte(data), &user)
			if e != nil {
				c.logger.Warn("failed to unmarshal data", zap.Error(e))
				errs = append(errs, e)
				return
			}
			users = append(users, user)

		default:
			c.logger.Warn("invalid value type", zap.String("type", fmt.Sprintf("%T", val)))
			// TODO add warning log
		}
	}
	return //nolint: nakedret // TODO: add directive
}
