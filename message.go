// Package slack - message
package slack

import (
	"context"
	"encoding/json"
	"fmt"
)

type (
	message struct {
		// Channel channel, private group, or IM channel to send message to. Can be an encoded ID, or a name
		Channel string `json:"channel"`

		// Text main body text of the message. If blocks field is passed, text is used as a fallback string to display
		// in notifications
		Text string `json:"text"`

		//AsUser post the message as the authed user, instead of as a bot
		AsUser *bool `json:"as_user,omitempty"`

		// Attachments to the message
		Attachments []Attachment `json:"attachments,omitempty"`

		// IconEmoji emoji to use as the icon for this message. Overrides icon_url
		IconEmoji string `json:"icon_emoji,omitempty"`

		// IconUrl url to an image to use as the icon for this message. Must be used in conjunction with as_user set to
		// false, otherwise ignored
		IconUrl string `json:"icon_url,omitempty"`

		// LinkNames find and link channel names and usernames.
		LinkNames *bool `json:"link_names,omitempty"`

		// Markdown disable Slack markup parsing by setting to false
		Markdown *bool `json:"mrkdwn,omitempty"`

		// ThreadTimestamp provide another message's ts value to make this message a repl
		ThreadTimestamp string `json:"thread_ts,omitempty"`

		// UnfurlLinks pass true to enable unfurling of primarily text-based content
		UnfurlLinks *bool `json:"unfurl_links,omitempty"`

		// UnfurlMedia pass false to disable unfurling of media content.
		UnfurlMedia *bool `json:"unfurl_media,omitempty"`

		// Username set bot's user name. Must be used in conjunction with as_user set to false, otherwise ignored.
		Username string `json:"username,omitempty"`
	}

	Attachment struct {
		// Fallback a plain text summary of the attachment used in clients that don't show formatted text
		Fallback string `json:"fallback,omitempty"`

		// Color changes the color of the border on the left side of this attachment from the default gray.
		Color string `json:"color,omitempty"`

		// Pretext text that appears above the message attachment block
		Pretext string `json:"pretext,omitempty"`

		// AuthorName small text used to display the author's name.
		AuthorName string `json:"author_name,omitempty"`

		// AuthorLink a valid URL that will hyperlink the author_name text. Will only work if author_name is present.
		AuthorLink string `json:"author_link,omitempty"`

		// AuthorIcon a valid URL that displays a small 16px by 16px image to the left of the author_name text. Will
		// only work if author_name is present.
		AuthorIcon string `json:"author_icon,omitempty"`

		// Title large title text near the top of the attachment.
		Title string `json:"title,omitempty"`

		// TitleLink a valid URL that turns the title text into a hyperlink.
		TitleLink string `json:"title_link,omitempty"`

		// The main body text of the attachment. It can be formatted as plain text, or with mrkdwn by including it in
		// the mrkdwn_in field. The content will automatically collapse if it contains 700+ characters or 5+ linebreaks,
		// and will display a "Show more..." link to expand the content.
		Text string `json:"text,omitempty"`

		// ImageUrl a valid URL to an image file that will be displayed at the bottom of the attachment. Supports GIF,
		// JPEG, PNG, and BMP formats.
		//
		// Large images will be resized to a maximum width of 360px or a maximum height of
		// 500px, while still maintaining the original aspect ratio. Cannot be used with thumb_url.
		ImageUrl string `json:"image_url,omitempty"`

		// ThumbUrl a valid URL to an image file that will be displayed as a thumbnail on the right side of a message
		// attachment. Currently supports the following formats: GIF, JPEG, PNG, and BMP.

		// The thumbnail's longest dimension will be scaled down to 75px while maintaining the aspect ratio of the image.
		// The filesize of the image must also be less than 500 KB.
		//
		// For best results, please use images that are already 75px by 75px.
		ThumbUrl string `json:"thumb_url,omitempty"`

		// Footer some brief text to help contextualize and identify an attachment. Limited to 300 characters, and may
		// be truncated further when displayed to users in environments with limited screen real estate.
		Footer string `json:"footer,omitempty"`

		// FooterIcon a valid URL to an image file that will be displayed beside the footer text. Will only work if
		// author_name is present. Will be rendered  at 16px by 16px. It's best to use an image that is similarly sized.
		FooterIcon string `json:"footer_icon,omitempty"`
	}
)

// MessagePosted entity
type MessagePosted struct {
	// Ok indicates success or failure
	Ok bool `json:"ok"`

	// Error short machine-readable error code
	Error string `json:"error"`

	// Channel where message was posted
	Channel string `json:"channel"`

	// Timestamp can be used in other message to reply
	Timestamp string `json:"timestamp"`

	// Message message body
	Message struct {
		// Text message text
		Text string `json:"text"`

		// Username bot's username
		Username string `json:"username"`

		// BotID id of bot
		BotID string `json:"bot_id"`

		// Attachments list
		Attachments []struct {
			// Text attachment's text
			Text string `json:"text"`

			// ID attachment's id
			ID int `json:"id"`

			// Fallback attachment's fallback
			Fallback string `json:"fallback"`
		} `json:"attachments"`

		// Type always message
		Type string `json:"type"`

		// SubType sub type of message
		SubType string `json:"subtype"`

		// Timestamp can be used in other message to reply
		Timestamp string `json:"timestamp"`
	} `json:"message"`
}

func postMessage(ctx context.Context, c *client, text, channel string, opts ...MsgOption) (MessagePosted, error) {
	message := message{
		Text:    text,
		Channel: channel,
	}

	for _, opt := range opts {
		opt.apply(&message)
	}

	data, err := json.Marshal(message)
	if err != nil {
		return MessagePosted{}, fmt.Errorf("can't marshal request message: %w", err)
	}

	var resp []byte
	if resp, err = c.post(ctx, "chat.postMessage", data); err != nil {
		return MessagePosted{}, err
	}

	var posted MessagePosted
	if err = json.Unmarshal(resp, &posted); err != nil {
		return MessagePosted{}, fmt.Errorf("can't unmarshal posted message data: %w", err)
	}

	return posted, nil
}
