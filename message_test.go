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

func TestMockClient_PostMessage(t *testing.T) {
	t.Run("positive case", func(t *testing.T) {
		var (
			baseUrl = "http://test.slack.com/api"
			message = "test_message"
			channel = "test_channel"
		)

		httpClient := new(slack.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
			req, ok := args.Get(0).(*http.Request)

			assert.True(t, ok)
			assert.Equal(t, baseUrl+"/"+"chat.postMessage", req.URL.String())

			expRequest := fmt.Sprintf("{\"channel\":\"%s\",\"text\":\"%s\"}", channel, message)
			request, err := ioutil.ReadAll(req.Body)
			assert.NoError(t, err)

			assert.Equal(t, expRequest, string(request))

		}).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(fmt.Sprintf("{\"channel\":\"%s\"}", channel)))),
			StatusCode: http.StatusOK,
		}, nil)

		c := slack.NewClient("test_token", slack.WithBaseUrl(baseUrl), slack.WithHttpClient(httpClient))

		resp, err := c.PostMessage(context.Background(), message, channel)
		assert.NoError(t, err)
		assert.Equal(t, slack.MessagePosted{Channel: channel}, resp)
	})

	t.Run("error on sending Message", func(t *testing.T) {
		var (
			expErr  = errors.New("test error")
			message = "test_message"
			channel = "test_channel"
		)

		httpClient := new(slack.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, expErr)

		client := slack.NewClient(
			"test_token",
			slack.WithHttpClient(httpClient),
		)

		resp, err := client.PostMessage(context.Background(), message, channel)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expErr))
		assert.Equal(t, slack.MessagePosted{}, resp)
	})

	t.Run("error on unmarshal response", func(t *testing.T) {
		var (
			message = "test_message"
			channel = "test_channel"
		)

		httpClient := new(slack.MockHTTPClient)
		httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("{"))),
			StatusCode: http.StatusOK,
		}, nil)

		client := slack.NewClient(
			"test_token",
			slack.WithHttpClient(httpClient),
		)

		discussion, err := client.PostMessage(context.Background(), message, channel)
		assert.Error(t, err)
		assert.Equal(t, slack.MessagePosted{}, discussion)
	})
}
