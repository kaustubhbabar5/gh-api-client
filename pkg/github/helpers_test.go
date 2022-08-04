package github //nolint: testpackage // TODO: add directive

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func (s *GithubTestSuite) TestNewRequest() {
	testCases := []struct {
		testName  string
		method    string
		path      string
		body      any
		errString string
	}{
		{
			testName:  "test01",
			method:    "GET",
			path:      "",
			body:      nil,
			errString: "",
		},
		{
			testName:  "test02",
			method:    "GET",
			path:      ":",
			body:      nil,
			errString: "parse \":\": missing protocol scheme",
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.testName, func() {
			_, err := s.ghClient.newRequest(testCase.method, testCase.path, testCase.body)
			if testCase.errString != "" {
				s.NotNil(err)
				s.EqualError(err, testCase.errString)
				return
			}

			s.Nil(err)
		})
	}
}

func (s *GithubTestSuite) TestValidateResponse() {
	testCases := []struct {
		testName  string
		res       *http.Response
		errString string
	}{
		{
			testName:  "test01",
			res:       &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(``))},
			errString: "",
		},
		{
			testName:  "test01",
			res:       &http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(strings.NewReader(``))},
			errString: "http request failed with status_code:500 and body :",
		},
		{
			testName:  "test01",
			res:       &http.Response{StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(strings.NewReader(``))},
			errString: "http request failed with status_code:500 and body :",
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.testName, func() {
			_, err := validateResponse(testCase.res)
			if testCase.errString != "" {
				s.NotNil(err)
				return
			}

			s.Nil(err)
		})
	}
}
