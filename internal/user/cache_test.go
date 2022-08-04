package user //nolint:testpackage // TODO add directive

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/kaustubhbabar5/gh-api-client/adapters/cache"
	"github.com/kaustubhbabar5/gh-api-client/pkg/config"
	"github.com/kaustubhbabar5/gh-api-client/pkg/github"
	"github.com/kaustubhbabar5/gh-api-client/pkg/logger"
	"github.com/stretchr/testify/suite"
)

type CacheTestSuite struct {
	suite.Suite
	c            *c
	testUser01   string
	testUser02   string
	testUser03   string
	testUserInfo string
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (s *CacheTestSuite) SetupSuite() {
	logger, err := logger.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
	}

	// TODO: use different config for testing
	config, err := config.Load("../..", logger)
	if err != nil {
		logger.Sugar().Error(err)
		return
	}

	redisClient, err := cache.NewRedisClient(config.RedisURL, config.RedisPasswordKey)
	if err != nil {
		logger.Sugar().Error(err)
		return
	}

	s.c = NewCache(logger, redisClient)

	s.testUser01 = "testUser1"
	s.testUser02 = "testUser3"
	s.testUser03 = "testUser2"

	s.testUserInfo = "{\"name\":\"Kaustubh Babar\",\"login\":\"kaustubhbabar5\",\"followers\":2,\"public_repos\":6}"
}

func (s *CacheTestSuite) BeforeTest(suiteName, testName string) {
	switch testName { //nolint:gocritic // TODO: directive
	case "TestGetUsers":
		err := s.c.rc.MSet(
			context.Background(),
			userCachePrefix+s.testUser01, s.testUserInfo,
			userCachePrefix+s.testUser02, s.testUserInfo,
			userCachePrefix+s.testUser03, s.testUserInfo,
		).Err()
		s.Nil(err)

		// default:
	}
}

// The TearDownTest method will be run after every test in the suite.
func (s *CacheTestSuite) TearDownTest() {
	s.c.rc.Del(
		context.Background(),
		userCachePrefix+s.testUser01,
		userCachePrefix+s.testUser02,
		userCachePrefix+s.testUser03,
	)
}
func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}

func (s *CacheTestSuite) TestCacheUserInfo() {
	testCases := []struct {
		testName  string
		usernames []github.User
		duration  time.Duration
	}{
		{
			testName: "test01",
			usernames: []github.User{
				{
					Name:        "Kaustubh Babar",
					Login:       "kaustubhbabar5",
					Company:     "",
					Followers:   1,
					PublicRepos: 2,
				},
				{
					Name:        "Kaustubh Babar",
					Login:       "kaustubhbabar52",
					Company:     "",
					Followers:   1,
					PublicRepos: 2,
				},
			},
			duration: 5 * time.Second,
		},
	}
	for _, testCase := range testCases {
		s.Run(testCase.testName, func() {
			err := s.c.CacheUserInfo(testCase.usernames, testCase.duration)
			s.Empty(err)

			time.Sleep(testCase.duration + time.Second)

			_, notFoundUsers, _ := s.c.GetUserInfo([]string{testCase.usernames[0].Login})
			s.NotEmpty(notFoundUsers)
		})
	}
}

func (s *CacheTestSuite) TestGetUsers() {
	testCases := []struct {
		testName           string
		usernames          []string
		expectedUsersCount int
	}{
		{
			testName:           "test1",
			usernames:          []string{s.testUser01, s.testUser02, s.testUser03},
			expectedUsersCount: 3,
		},
		{
			testName:           "test2",
			usernames:          []string{s.testUser01, s.testUser02, s.testUser03, "invalid_user"},
			expectedUsersCount: 3,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.testName, func() {
			users, _, errs := s.c.GetUserInfo(testCase.usernames)
			s.Empty(errs)
			s.Equal(testCase.expectedUsersCount, len(users))
		})
	}
}
