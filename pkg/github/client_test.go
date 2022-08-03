package github

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
	}

	// TODO: use different config for testing
	config, err := config.Load("../..", logger)
	if err != nil {
		logger.Sugar().Error(err)
	}

	httpClient := &http.Client{
		Timeout: time.Duration(config.ReadTimeout) * time.Second,
	}

	s.ghClient = NewClient(httpClient, config.GithubAuthToken)
}

func TestGithubTestSuite(t *testing.T) {
	suite.Run(t, new(GithubTestSuite))
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
				PublicRepos: 6,
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
			s.Equal(true, reflect.DeepEqual(user, testCase.expectedUser), fmt.Sprintf("expected %v, got %v", testCase.expectedUser, user))

		})
	}
}
