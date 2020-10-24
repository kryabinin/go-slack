// Package slack - client
package slack

//go:generate mockery -case=underscore -inpkg -name=Client
//go:generate mockery -case=underscore -inpkg -name=HTTPClient

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const defaultBaseUrl = "https://slack.com/api/"

// HTTPClient interface to replace default http client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type (
	// Client provides api to work with slack entities
	Client interface {
		// PostMessage send Message to a channel
		PostMessage(ctx context.Context, message string, channel string, opts ...MsgOption) (MessagePosted, error)

		// GetUserByEmail find a user with an email address.
		GetUserByEmail(ctx context.Context, email string) (User, error)

		// SendRequest send http request to slack
		SendRequest(ctx context.Context, method string, path string, data []byte) ([]byte, error)
	}

	client struct {
		token      string
		baseUrl    string
		httpClient HTTPClient
	}
)

// NewClient is client constructor
func NewClient(token string, opts ...ClientOption) Client {
	c := &client{
		token:      token,
		baseUrl:    defaultBaseUrl,
		httpClient: &http.Client{},
	}

	for _, opt := range opts {
		opt.apply(c)
	}

	return c
}

// GetUserByEmail implementation
func (c *client) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return getUserByEmail(ctx, c, email)
}

// PostMessage implementation
func (c *client) PostMessage(ctx context.Context, text, channel string, opts ...MsgOption) (MessagePosted, error) {
	return postMessage(ctx, c, text, channel, opts...)
}

func (c *client) get(ctx context.Context, path string) ([]byte, error) {
	return c.SendRequest(ctx, http.MethodGet, path, nil)
}

func (c *client) post(ctx context.Context, path string, data []byte) ([]byte, error) {
	return c.SendRequest(ctx, http.MethodPost, path, data)
}

// SendRequest implementation
func (c *client) SendRequest(ctx context.Context, method string, path string, data []byte) ([]byte, error) {
	if len(c.token) == 0 {
		return nil, errors.New("token is required")
	}

	req, err := http.NewRequest(method, c.baseUrl+"/"+path, bytes.NewReader(data))
	if nil != err {
		return nil, fmt.Errorf("can't create http request: %w", err)
	}

	req = req.WithContext(ctx)

	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	var resp *http.Response
	if resp, err = c.httpClient.Do(req); nil != err {
		return nil, fmt.Errorf("can't send http request: %w", err)
	}

	if http.StatusOK != resp.StatusCode {
		return nil, fmt.Errorf("slack respond with %d status code", resp.StatusCode)
	}

	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); nil != err {
		return nil, fmt.Errorf("can't read response body: %w", err)
	}

	return body, nil
}
