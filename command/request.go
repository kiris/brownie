package command

import (
	"github.com/nlopes/slack"
)

type Request struct {
	listener    *Listener
	Event       *slack.MessageEvent
	Mention     string
	CommandName string
	CommandArgs []string
}
