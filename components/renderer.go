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
	channel, ts, threadTs, attachments := component.Render()
	attachmentOption := slack.MsgOptionAttachments(attachments ...)

	if ts != "" {
		if _, _, _, err := r.Client.UpdateMessage(channel, ts, attachmentOption); err != nil {
			return errors.Wrapf(err, "failed to render with update: Channel = %s, Ts = %s", channel, ts)
		}
	} else {
		options := []slack.MsgOption {
			slack.MsgOptionAttachments(attachments ...),
		}
		if threadTs != "" {
			options = append(options, slack.MsgOptionTS(threadTs))
		}
		if _, _, err := r.Client.PostMessage(channel, options ...); err != nil {
			return errors.Wrapf(err, "failed to render: Channel = %s", channel)
		}
	}

	return nil
}
