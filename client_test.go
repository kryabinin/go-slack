package slack_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kryabinin/go-slack"
)

func TestClient_SendRequest(t *testing.T) {
	t.Run("positive case", func(t *testing.T) {
		var (
			token       = "test_token"
			baseUrl     = "http://test.slack.com/api"
			path        = "/test/path"
			method      = http.MethodGet
			expResponse = []byte(`{"test": "passed"}`)
		)

		httpClient := new(slack.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
			req, ok := args.Get(0).(*http.Request)
			assert.True(t, ok)

			assert.Equal(t, method, req.Method)
			assert.Equal(t, baseUrl+"/"+path, req.URL.String())

			assert.Equal(t, "Bearer "+token, req.Header.Get("Authorization"))
			assert.Equal(t, "application/json; charset=utf-8", req.Header.Get("Content-Type"))
		}).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader(expResponse)),
			StatusCode: http.StatusOK,
		}, nil)

		client := slack.NewClient(
			token,
			slack.WithBaseUrl(baseUrl),
			slack.WithHttpClient(httpClient),
		)

		resp, err := client.SendRequest(context.Background(), method, path, nil)
		assert.NoError(t, err)
		assert.Equal(t, expResponse, resp)
	})

	t.Run("error on sending http request", func(t *testing.T) {
		expErr := errors.New("test error")

		httpClient := new(slack.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, expErr)

		client := slack.NewClient(
			"test_token",
			slack.WithHttpClient(httpClient),
		)

		resp, err := client.SendRequest(context.Background(), http.MethodGet, "/test/path", nil)
		assert.Equal(t, []byte(nil), resp)
		assert.True(t, errors.Is(err, expErr))
	})

	t.Run("non 200 status", func(t *testing.T) {
		expStatus := http.StatusInternalServerError

		httpClient := new(slack.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
			StatusCode: expStatus,
		}, nil)

		client := slack.NewClient(
			"test_token",
			slack.WithHttpClient(httpClient),
		)

		resp, err := client.SendRequest(context.Background(), http.MethodGet, "/test/path", nil)
		assert.Equal(t, []byte(nil), resp)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("slack respond with %d status code", expStatus))
	})

	t.Run("error on read response body", func(t *testing.T) {
		expStatus := http.StatusInternalServerError

		httpClient := new(slack.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
			StatusCode: expStatus,
		}, nil)

		client := slack.NewClient(
			"test_token",
			slack.WithHttpClient(httpClient),
		)

		resp, err := client.SendRequest(context.Background(), http.MethodGet, "/test/path", nil)
		assert.Equal(t, []byte(nil), resp)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), fmt.Sprintf("slack respond with %d status code", expStatus))
	})
}
