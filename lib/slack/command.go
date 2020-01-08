package slack

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)


type CommandListener struct {
	Rtm      *slack.RTM
	botId    string
	botName  string
	handlers map[string]CommandHandler
}


type CommandHandler interface {
	ExecCommand(req *CommandRequest) error
}

type CommandRequest struct {
	listener    *CommandListener
	Event       *slack.MessageEvent
	Mention     string
	CommandName string
	CommandArgs []string
}

func NewCommandListener(client *slack.Client) *CommandListener {
	return &CommandListener{
		Rtm     : client.NewRTM(),
		handlers: map[string]CommandHandler{},
	}
}

func (l *CommandListener) ListenAndResponse() error {
	go l.Rtm.ManageConnection()

	for msg := range l.Rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectingEvent:
			log.Info("Slack RTM connecting.")

		case *slack.ConnectedEvent:
			log.WithFields(log.Fields{
				"id":   ev.Info.User.ID,
				"CommandName": ev.Info.User.Name,
			}).Info("Slack RTM connected.")

			l.botId = ev.Info.User.ID
			l.botName = ev.Info.User.Name

		case *slack.MessageEvent:
			l.handleMessageEvent(ev)

		case *slack.InvalidAuthEvent:
			return errors.New("invalid credentials")

		default:
			// Ignore other events.
			// fmt.Printf("other: %v\n", ev)
		}
	}

	return nil
}

func (l *CommandListener) Handle(name string, handler CommandHandler) {
	l.handlers[name] = handler
}

func (l *CommandListener) handleMessageEvent(ev *slack.MessageEvent) {
	req := l.parseRequest(ev)
	if req == nil {
		return
	}
	// Only response Mention mentionTo bot. Ignore else.
	if !l.mentionToMe(req) {
		return
	}

	if l.messageFromBot(ev) {
		log.Debug("ignore bot message.")
		return
	}

	fmt.Println(req.CommandName)
	handler, ok := l.handlers[req.CommandName]
	if !ok {
		log.WithField("command", req.CommandName).Warn("unknown command")
		return
	}

	if err := handler.ExecCommand(req); err != nil {
		log.WithError(err).WithField("command", req.CommandName).Error("failed exec command")
		return
	}
}

func (l *CommandListener) parseRequest(event *slack.MessageEvent) *CommandRequest {
	tokens := strings.Split(strings.TrimSpace(event.Msg.Text), " ")
	if len(tokens) == 1 {
		return nil
	}

	return &CommandRequest{
		listener   : l,
		Event      : event,
		Mention    : tokens[0],
		CommandName: l.parseCommandName(tokens),
		CommandArgs: l.parseCommandArgs(tokens),
	}
}

func (l *CommandListener) parseCommandName(tokens []string) string {
	if len(tokens) < 2 {
		return "help"
	} else {
		return tokens[1]
	}
}

func (l *CommandListener) parseCommandArgs(tokens []string) []string {
	return tokens[2:]
}

func (l *CommandListener) messageTokens(ev *slack.MessageEvent) []string {
	return strings.Split(strings.TrimSpace(ev.Msg.Text), " ")
}

func (l *CommandListener) messageFromBot(ev *slack.MessageEvent) bool {
	return ev.Msg.SubType == "bot_message"
}


func (l *CommandListener) mentionToMe(req *CommandRequest) bool {
	return req.Mention == fmt.Sprintf("<@%s>", l.botId)
}