package user

import (
	"time"

	"github.com/kaustubhbabar5/gh-api-client/pkg/github"
	"go.uber.org/zap"
)

const (
	cacheDuration = 2 // in minutes
)

type Service interface {
	// fetches users information by username
	GetUsers(usernames []string) (users github.Users, notFoundUsers []string, errs []error)

	// fetches users information by username, sorted alphabetically by login case insensitive
	GetUsersSorted(usernames []string) (
		users github.Users,
		notFoundUsers []string,
		errs []error,
	)
}

type service struct {
	logger       *zap.Logger
	userCache    Cache
	githubClient github.Client
}

func NewService(logger *zap.Logger, userCache Cache, githubClient github.Client) *service {
	return &service{logger, userCache, githubClient}
}

// GetUsers retrieves users from the cache,
// for the user info that is not cached it retrieves information from github API
// returns back with  users, users that were not found on github and errors if any.
func (s *service) GetUsers(usernames []string) ( //nolint: nonamedreturns // TODO: add directive
	users github.Users,
	notFoundUsers []string,
	errs []error,
) {
	cachedUsers, notCachedUsers, errs := s.userCache.GetUserInfo(usernames)
	if len(errs) != 0 {
		s.logger.Warn("failed to get user info from cache %w", zap.Errors("errors", errs))
	}

	if len(cachedUsers) == len(usernames) {
		users = cachedUsers
		return
	}

	users, notFoundUsers, errs = s.githubClient.GetUsers(notCachedUsers)

	if len(users) == 0 {
		users = cachedUsers
		return
	}

	users = append(users, cachedUsers...)

	err := s.userCache.CacheUserInfo(users, cacheDuration*time.Minute)
	if len(err) != 0 {
		s.logger.Warn("failed to cache User info", zap.Errors("errors", err))
	}

	return //nolint: nakedret // TODO: add directive
}

func (s *service) GetUsersSorted(usernames []string) ( //nolint: nonamedreturns // TODO: add directive
	users github.Users,
	notFoundUsers []string,
	errs []error,
) {
	users, notFoundUsers, errs = s.GetUsers(usernames)

	users.SortByLogin()

	return //nolint: nakedret // TODO: add directive
}
