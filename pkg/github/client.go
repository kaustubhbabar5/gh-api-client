package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	cerrors "github.com/kaustubhbabar5/gh-api-client/pkg/errors"
)

const (
	BASE_URL_STRING = "https://api.github.com/"
)

type Client interface {
	GetUser(username string) (User, error)
}

type client struct {
	httpClient *http.Client
	baseURL    *url.URL
	authToken  string
}

// creates a new github client with given http.Client and config.
// if http Client is nil it will create a new http.Client
func NewClient(httpClient *http.Client, authToken string) *client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL, _ := url.Parse(BASE_URL_STRING)
	authToken = fmt.Sprintf("token %s", authToken)

	return &client{
		httpClient,
		baseURL,
		authToken,
	}
}

// Fetches information about user by username,
// will return back with custom error `Not Found` if the user is not found
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
