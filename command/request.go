package command

import (
	"fmt"
	"github.com/nlopes/slack"
)

type Request struct {
	listener    *Listener
	Event       *slack.MessageEvent
	Mention     string
	CommandName string
	CommandArgs []string
}

func (r *Request) ResponseAttachmentsMessage(attachments ...slack.Attachment) (string, error) {
	return r.ResponseMessage(slack.MsgOptionAttachments(attachments ...))
}

func (r *Request) ResponseMessage(options ...slack.MsgOption) (string, error) {
	if _, ts, err := r.listener.Rtm.PostMessage(r.Event.Channel, options ...); err != nil {
		return "", fmt.Errorf("failed to post message: %s", err)
	} else {
		return ts, nil
	}
}
