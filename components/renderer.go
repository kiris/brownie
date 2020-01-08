package components

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"net/http"
)

type Component interface {
	Render() (string, string, string, []slack.Attachment)
}

type Renderer interface {
	Render(component Component) error
}

type InteractionRenderer struct {
	Writer http.ResponseWriter
	Callback *slack.InteractionCallback
}

func (r InteractionRenderer) Render(component Component) error {
	message := r.Callback.OriginalMessage
	_, ts, _, attachments := component.Render()
	if ts != "" {
		message.ReplaceOriginal = true
	}

	message.Attachments = attachments
	r.Writer.Header().Add("Content-type", "application/json")
	r.Writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(r.Writer).Encode(&message); err != nil {
		return errors.Wrapf(err, "failed to render with interaction")
	}

	return nil
}

type ApiRenderer struct {
	Client *slack.Client

}

func (r *ApiRenderer) Render(component Component) error {
	channel, ts, threadTS, attachments := component.Render()
	attachmentOption := slack.MsgOptionAttachments(attachments ...)

	options := []slack.MsgOption { attachmentOption }
	if ts != "" {
		options = append(options, slack.MsgOptionUpdate(ts))
	}
	if threadTS != "" {
		options = append(options, slack.MsgOptionTS(threadTS))
	}
	if _, _, err := r.Client.PostMessage(channel, options ...); err != nil {
		return errors.Wrapf(err, "failed to render: channel = %s, ts = %s, threadTS = %s", channel, ts, threadTS)
	}

	return nil
}
