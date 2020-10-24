// Package slack - Message options
package slack

type (
	// MsgOption to apply optional options to Message
	MsgOption interface {
		apply(msg *Message)
	}

	asUser struct {
		val bool
	}

	addAttachment struct {
		attachment Attachment
	}
)

// AsUser pass true to post the Message as the authed user, instead of as a bot
func AsUser(val bool) MsgOption {
	return &asUser{val: val}
}

func (opt *asUser) apply(msg *Message) {
	msg.AsUser = &opt.val
}

// AddAttachment adds attachment to the Message
func AddAttachment(attachment Attachment) MsgOption {
	return &addAttachment{attachment: attachment}
}

func (opt *addAttachment) apply(msg *Message) {
	msg.Attachments = append(msg.Attachments, opt.attachment)
}
