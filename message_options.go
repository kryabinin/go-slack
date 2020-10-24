// Package slack - message options
package slack

type (
	// MsgOption to apply optional options to message
	MsgOption interface {
		apply(msg *message)
	}

	asUser struct {
		val bool
	}

	addAttachment struct {
		attachment Attachment
	}
)

// AsUser pass true to post the message as the authed user, instead of as a bot
func AsUser(val bool) MsgOption {
	return &asUser{val: val}
}

func (opt *asUser) apply(msg *message) {
	msg.AsUser = &opt.val
}

// AddAttachment adds attachment to the message
func AddAttachment(attachment Attachment) MsgOption {
	return &addAttachment{attachment: attachment}
}

func (opt *addAttachment) apply(msg *message) {
	msg.Attachments = append(msg.Attachments, opt.attachment)
}
