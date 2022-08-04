package github //nolint:testpackage // TODO add directive

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/kaustubhbabar5/gh-api-client/pkg/config"
	"github.com/kaustubhbabar5/gh-api-client/pkg/logger"
	"github.com/stretchr/testify/suite"
)

type GithubTestSuite struct {
	suite.Suite
	ghClient *client
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (s *GithubTestSuite) SetupSuite() {
	logger, err := logger.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
		return
	}

	// TODO: use different config for testing
	config, err := config.Load("../..", logger)
	if err != nil {
		logger.Sugar().Error(err)
		return
	}
	httpClient := &http.Client{
		Timeout: time.Duration(config.ReadTimeout) * time.Second,
	}

	s.ghClient = NewClient(httpClient, config.GithubAuthTokenKey)
}

func TestGithubTestSuite(t *testing.T) {
	suite.Run(t, new(GithubTestSuite))
}

func (s *GithubTestSuite) TestNewClient() {
	logger, err := logger.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
	}

	// TODO: use different config for testing
	config, err := config.Load("../..", logger)
	if err != nil {
		logger.Sugar().Error(err)
	}

	httpClient := &http.Client{
		Timeout: time.Duration(config.ReadTimeout) * time.Second,
	}

	githubClient := NewClient(httpClient, config.GithubAuthTokenKey)

	s.Equal(BaseURLString, githubClient.baseURL.String())

	githubClient2 := NewClient(httpClient, config.GithubAuthTokenKey)

	if githubClient == githubClient2 {
		s.Fail("both clients should differ, but they are same")
	}
}

func (s *GithubTestSuite) TestGetUser() {
	testCases := []struct {
		testName     string
		username     string
		expectedUser User
		errString    string
	}{
		{
			testName: "test01",
			username: "kaustubhbabar5",
			expectedUser: User{
				Name:        "Kaustubh Babar",
				Login:       "kaustubhbabar5",
				Company:     "",
				Followers:   2,
				PublicRepos: 7,
			},
			errString: "",
		},
		{
			testName:     "test02",
			username:     "kaustubhbabar5_invalid",
			expectedUser: User{},
			errString:    "user not found",
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.testName, func() {
			user, err := s.ghClient.GetUser(testCase.username)
			if testCase.errString != "" {
				s.NotNil(err)
				s.EqualError(err, testCase.errString)
				return
			}
			s.Nil(err)
			s.Equal(
				true,
				reflect.DeepEqual(
					user,
					testCase.expectedUser,
				),
				fmt.Sprintf("expected %v, got %v", testCase.expectedUser, user),
			)
		})
	}
}

func (s *GithubTestSuite) TestGetUsers() {
	testCases := []struct {
		testName      string
		usernames     []string
		userCount     int
		errCount      int
		notFoundCount int
	}{
		{
			usernames: []string{
				"kaustubhbabar5",
				"dhruvikn",
				"exagil",
				"GotamDahiya",
				"manya28",
				"Mayurkumawat012",
				"Shubham6147",
				"suryakanigolla",
				"vgnh",
				"nikhil-kawa",
			},
			userCount:     10,
			errCount:      0,
			notFoundCount: 0,
		},
		{
			usernames: []string{
				"kaustubhbabar5",
				"dhruvikn",
				"exagil--invalid",
				"GotamDahiya--invalid",
				"manya28",
				"Mayurkumawat012--invalid",
				"Shubham6147",
				"suryakanigolla",
				"vgnh",
				"nikhil-kawa-invalid",
			},
			userCount:     6,
			errCount:      0,
			notFoundCount: 4,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.testName, func() {
			user, usernamesNotFound, errs := s.ghClient.GetUsers(testCase.usernames)
			s.Equal(testCase.errCount, len(errs))
			s.Equal(testCase.notFoundCount, len(usernamesNotFound))
			s.Equal(testCase.userCount, len(user))
		})
	}
}
