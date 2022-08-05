package user //nolint:testpackage // TODO add directive

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/kaustubhbabar5/gh-api-client/internal/user/mocks"
	"github.com/kaustubhbabar5/gh-api-client/pkg/github"
	githubmocks "github.com/kaustubhbabar5/gh-api-client/pkg/github/mocks"

	"github.com/kaustubhbabar5/gh-api-client/pkg/logger"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	service          *service
	userCacheMock    *mocks.Cache
	gitHubClientMock *githubmocks.Client
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (s *ServiceTestSuite) SetupSuite() {
	logger, err := logger.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
	}

	s.userCacheMock = &mocks.Cache{}
	s.gitHubClientMock = &githubmocks.Client{}

	s.service = NewService(logger, s.userCacheMock, s.gitHubClientMock)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestGetUsers() {
	usernames := []string{"one", "two"}
	expectedUsers := github.Users{
		{
			Name:        "Kaustubh Babar",
			Login:       "one",
			Company:     "",
			Followers:   1,
			PublicRepos: 2,
		},
		{
			Name:        "Kaustubh Babar",
			Login:       "two",
			Company:     "",
			Followers:   1,
			PublicRepos: 2,
		},
	}

	s.gitHubClientMock.On("GetUsers", usernames).Return(expectedUsers, []string{}, nil)

	s.userCacheMock.On("GetUserInfo", usernames).Return(github.Users{}, usernames, nil)
	s.userCacheMock.On("CacheUserInfo", expectedUsers, 2*time.Minute).Return(nil)

	actualUsers, notFoundUsers, errs := s.service.GetUsers(usernames)
	s.Empty(errs)
	s.Empty(notFoundUsers)
	s.Equal(true, reflect.DeepEqual(expectedUsers, actualUsers))

	s.gitHubClientMock.AssertExpectations(s.T())
}

func (s *ServiceTestSuite) TestGetUsersFail() {
	usernames := []string{"three", "four"}

	s.userCacheMock.On("GetUserInfo", usernames).Return(github.Users{}, usernames, nil)
	s.gitHubClientMock.On("GetUsers", usernames).Return(github.Users{}, usernames, nil)

	actualUsers, notFoundUsers, errs := s.service.GetUsers(usernames)
	s.Empty(errs)
	s.Empty(actualUsers)
	s.Equal(true, reflect.DeepEqual(usernames, notFoundUsers))
}
