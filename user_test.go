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

func TestClient_GetUserByEmail(t *testing.T) {
	t.Run("positive case", func(t *testing.T) {
		var (
			email   = "user@mail.ru"
			baseUrl = "http://test.slack.com/api"
		)

		httpClient := new(slack.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
			req, ok := args.Get(0).(*http.Request)

			assert.True(t, ok)
			assert.Equal(t, baseUrl+"/"+fmt.Sprintf("/api/users.lookupByEmail?email=%s", email), req.URL.String())
		}).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(fmt.Sprintf("{" +
				"\"ok\":true," +
				"\"user\":{\"profile\": {\"email\":\"%s\"}}" +
				"}", email)))),
			StatusCode: http.StatusOK,
		}, nil)

		client := slack.NewClient(
			"test_token",
			slack.WithBaseUrl(baseUrl),
			slack.WithHttpClient(httpClient),
		)

		user, err := client.GetUserByEmail(context.Background(), email)
		assert.NoError(t, err)
		assert.Equal(t, email, user.Profile.Email)
	})

	t.Run("error on getting user", func(t *testing.T) {
		expErr := errors.New("test error")

		httpClient := new(slack.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, expErr)

		client := slack.NewClient(
			"test_token",
			slack.WithHttpClient(httpClient),
		)

		user, err := client.GetUserByEmail(context.Background(), "test@email.com")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expErr))
		assert.Equal(t, slack.User{}, user)
	})

	t.Run("error on unmarshal response", func(t *testing.T) {
		httpClient := new(slack.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("{"))),
			StatusCode: http.StatusOK,
		}, nil)

		client := slack.NewClient(
			"test_token",
			slack.WithHttpClient(httpClient),
		)

		user, err := client.GetUserByEmail(context.Background(), "test@email.com")
		assert.Error(t, err)
		assert.Equal(t, slack.User{}, user)
	})

	t.Run("slack respond with error", func(t *testing.T) {
		httpClient := new(slack.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("{\"ok\":false,\"error\":\"test\"}"))),
			StatusCode: http.StatusOK,
		}, nil)

		client := slack.NewClient(
			"test_token",
			slack.WithHttpClient(httpClient),
		)

		user, err := client.GetUserByEmail(context.Background(), "test@mail.com")
		assert.Error(t, err)
		assert.Equal(t, slack.User{}, user)
	})
}
