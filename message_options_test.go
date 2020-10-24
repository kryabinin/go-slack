package slack_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kryabinin/go-slack"
)

func TestAsUser(t *testing.T) {
	for _, asUser := range []bool{true, false} {
		t.Run(fmt.Sprintf("%t", asUser), func(t *testing.T) {
			var (
				message = "test_message"
				channel = "test_channel"
			)

			httpClient := new(slack.MockHTTPClient)
			httpClient.On("Do", mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
				req, ok := args.Get(0).(*http.Request)
				assert.True(t, ok)

				expRequest := fmt.Sprintf("{"+
					"\"channel\":\"%s\","+
					"\"text\":\"%s\","+
					"\"as_user\":%t}", channel, message, asUser)
				request, err := ioutil.ReadAll(req.Body)
				assert.NoError(t, err)

				assert.Equal(t, expRequest, string(request))

			}).Return(&http.Response{
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
				StatusCode: http.StatusOK,
			}, nil)

			c := slack.NewClient("test_token", slack.WithHttpClient(httpClient))

			_, err := c.PostMessage(context.Background(), message, channel, slack.AsUser(asUser))
			assert.NoError(t, err)
		})
	}
}

func TestAddAttachment(t *testing.T) {
	var (
		text    = "test_text"
		channel = "test_channel"
	)
	testCases := []struct {
		name        string
		attachments []slack.Attachment
		expRequest  string
	}{
		{
			name:        "fallback",
			attachments: []slack.Attachment{{Fallback: "test_fallback1"}, {Fallback: "test_fallback2"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"fallback\":\"test_fallback1\"},{\"fallback\":\"test_fallback2\"}]"+
				"}", channel, text),
		},
		{
			name:        "color",
			attachments: []slack.Attachment{{Color: "red"}, {Color: "green"}, {Color: "blue"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"color\":\"red\"},{\"color\":\"green\"},{\"color\":\"blue\"}]"+
				"}", channel, text),
		},
		{
			name:        "pretext",
			attachments: []slack.Attachment{{Pretext: "test_pretext1"}, {Pretext: "test_pretext2"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"pretext\":\"test_pretext1\"},{\"pretext\":\"test_pretext2\"}]"+
				"}", channel, text),
		},
		{
			name:        "author_name",
			attachments: []slack.Attachment{{AuthorName: "test_author_name1"}, {AuthorName: "test_author_name2"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"author_name\":\"test_author_name1\"},{\"author_name\":\"test_author_name2\"}]"+
				"}", channel, text),
		},
		{
			name:        "author_link",
			attachments: []slack.Attachment{{AuthorLink: "test_author_link1"}, {AuthorLink: "test_author_link2"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"author_link\":\"test_author_link1\"},{\"author_link\":\"test_author_link2\"}]"+
				"}", channel, text),
		},
		{
			name:        "author_icon",
			attachments: []slack.Attachment{{AuthorIcon: "test_author_icon1"}, {AuthorIcon: "test_author_icon2"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"author_icon\":\"test_author_icon1\"},{\"author_icon\":\"test_author_icon2\"}]"+
				"}", channel, text),
		},
		{
			name:        "title",
			attachments: []slack.Attachment{{Title: "test_title1"}, {Title: "test_title2"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"title\":\"test_title1\"},{\"title\":\"test_title2\"}]"+
				"}", channel, text),
		},
		{
			name:        "title_link",
			attachments: []slack.Attachment{{TitleLink: "test_title_link1"}, {TitleLink: "test_title_link2"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"title_link\":\"test_title_link1\"},{\"title_link\":\"test_title_link2\"}]"+
				"}", channel, text),
		},
		{
			name:        "text",
			attachments: []slack.Attachment{{Text: "test_text1"}, {Text: "test_text2"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"text\":\"test_text1\"},{\"text\":\"test_text2\"}]"+
				"}", channel, text),
		},
		{
			name:        "image_url",
			attachments: []slack.Attachment{{ImageUrl: "test_image_url1"}, {ImageUrl: "test_image_url2"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"image_url\":\"test_image_url1\"},{\"image_url\":\"test_image_url2\"}]"+
				"}", channel, text),
		},
		{
			name:        "thumb_url",
			attachments: []slack.Attachment{{ThumbUrl: "test_thumb_url1"}, {ThumbUrl: "test_thumb_url2"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"thumb_url\":\"test_thumb_url1\"},{\"thumb_url\":\"test_thumb_url2\"}]"+
				"}", channel, text),
		},
		{
			name:        "footer",
			attachments: []slack.Attachment{{Footer: "test_footer1"}, {Footer: "test_footer2"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"footer\":\"test_footer1\"},{\"footer\":\"test_footer2\"}]"+
				"}", channel, text),
		},
		{
			name:        "footer_icon",
			attachments: []slack.Attachment{{FooterIcon: "test_footer_icon1"}, {FooterIcon: "test_footer_icon2"}},
			expRequest: fmt.Sprintf("{"+
				"\"channel\":\"%s\","+
				"\"text\":\"%s\","+
				"\"attachments\":[{\"footer_icon\":\"test_footer_icon1\"},{\"footer_icon\":\"test_footer_icon2\"}]"+
				"}", channel, text),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			httpClient := new(slack.MockHTTPClient)
			httpClient.On("Do", mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
				req, ok := args.Get(0).(*http.Request)
				assert.True(t, ok)

				request, err := ioutil.ReadAll(req.Body)
				assert.NoError(t, err)

				assert.Equal(t, testCase.expRequest, string(request))

			}).Return(&http.Response{
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
				StatusCode: http.StatusOK,
			}, nil)

			c := slack.NewClient("test_token", slack.WithHttpClient(httpClient))

			attachmentOpts := make([]slack.MsgOption, 0, len(testCase.attachments))
			for _, attachment := range testCase.attachments {
				attachmentOpts = append(attachmentOpts, slack.AddAttachment(attachment))
			}

			_, err := c.PostMessage(context.Background(), text, channel, attachmentOpts...)
			assert.NoError(t, err)
		})
	}
}
