package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	cerrors "github.com/kaustubhbabar5/gh-api-client/pkg/errors"
)

const (
	BaseURLString = "https://api.github.com/"
)

type Client interface {
	GetUser(username string) (User, error)
	GetUsers(usernames []string) (users []User, usersNotFound []string, errs []error)
}

type client struct {
	httpClient *http.Client
	baseURL    *url.URL
	authToken  string
}

// creates a new github client with given http.Client and config.
// if http Client is nil it will create a new http.Client.
func NewClient(httpClient *http.Client, authTokenKey string) *client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(BaseURLString)
	authToken := fmt.Sprintf("token %v", os.Getenv(authTokenKey))

	return &client{
		httpClient,
		baseURL,
		authToken,
	}
}

// Fetches information about user by username,
// will return back with custom error `Not Found` if the user is not found.
func (c *client) GetUser(username string) (User, error) {
	if username == "" {
		return User{}, errors.New("username empty")
	}
	path := fmt.Sprintf("users/%s", username)

	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return User{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer res.Body.Close()

	body, err := validateResponse(res)
	if err != nil {
		return User{}, err
	}

	if res.StatusCode == http.StatusNotFound {
		return User{}, cerrors.NewNotFound("user", string(body))
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// GetUsers returns information of multiple users by their usernames.
func (c *client) GetUsers(usernames []string) ( //nolint:nonamedreturns // named returns serves as documentation here
	users []User,
	usersNotFound []string,
	errs []error,
) {
	userCount := len(usernames)

	users = make([]User, 0)
	usersNotFound = make([]string, 0)
	errs = make([]error, 0)

	userChan := make(chan User, userCount)
	userNotFoundChan := make(chan string, userCount)
	errChan := make(chan error, userCount)

	for _, username := range usernames {
		go func(username string, userChan chan User, userNotFoundChan chan string, errChan chan error) {
			user, err := c.GetUser(username)
			if err != nil {
				var notFoundErr *cerrors.NotFoundError
				if errors.As(err, &notFoundErr) {
					userNotFoundChan <- username
					return
				}

				errChan <- fmt.Errorf("failed to get information for username:%s error:%w", username, err)
				return
			}

			userChan <- user
		}(username, userChan, userNotFoundChan, errChan)
	}

	for i := 1; i <= userCount; i++ {
		select {
		case user := <-userChan:
			users = append(users, user)

		case err := <-errChan:
			errs = append(errs, err)

		case username := <-userNotFoundChan:
			usersNotFound = append(usersNotFound, username)
		}
	}

	return //nolint:nakedret // added directive to func def
}
