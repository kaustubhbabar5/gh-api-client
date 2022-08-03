package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func (c *client) newRequest(method, path string, body any) (*http.Request, error) {
	url, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		// `application/vnd.github+json` is recommended by github in API docs ref: https://docs.github.com/en/rest/users/users#get-a-user--parameters
		req.Header.Set("Content-Type", "application/vnd.github+json")
	}

	if c.authToken != "" {
		req.Header.Set("Authorization", c.authToken)
	}

	return req, nil
}

// validateResponse checks the http.Response for errors and reads the body, returns error if status code is anything other than 200 or fails to read the body
func validateResponse(res *http.Response) ([]byte, error) {
	//TODO: handle rate limits errors here

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		return nil, fmt.Errorf("http request failed with status_code:%d and body :%s", res.StatusCode, body)
	}
	return body, nil
}
