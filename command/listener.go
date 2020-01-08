package command

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)


type Listener struct {
	Rtm      *slack.RTM
	botId    string
	botName  string
	handlers map[string]Handler
}


type Handler interface {
	ExecCommand(req *Request) error
}

func CreateListener(client *slack.Client) *Listener {
	return &Listener{
		Rtm     : client.NewRTM(),
		handlers: map[string]Handler{},
	}
}

func (l *Listener) ListenAndResponse() error {
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

func (l *Listener) Handle(name string, handler Handler) {
	l.handlers[name] = handler
}

func (l *Listener) handleMessageEvent(ev *slack.MessageEvent) {
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

func (l *Listener) parseRequest(event *slack.MessageEvent) *Request {
	tokens := strings.Split(strings.TrimSpace(event.Msg.Text), " ")
	if len(tokens) == 1 {
		return nil
	}

	return &Request{
		listener   : l,
		Event      : event,
		Mention    : tokens[0],
		CommandName: l.parseCommandName(tokens),
		CommandArgs: l.parseCommandArgs(tokens),
	}
}

func (l *Listener) parseCommandName(tokens []string) string {
	if len(tokens) < 2 {
		return "help"
	} else {
		return tokens[1]
	}
}

func (l *Listener) parseCommandArgs(tokens []string) []string {
	return tokens[2:]
}

func (l *Listener) messageTokens(ev *slack.MessageEvent) []string {
	return strings.Split(strings.TrimSpace(ev.Msg.Text), " ")
}

func (l *Listener) messageFromBot(ev *slack.MessageEvent) bool {
	return ev.Msg.SubType == "bot_message"
}


func (l *Listener) mentionToMe(req *Request) bool {
	return req.Mention == fmt.Sprintf("<@%s>", l.botId)
}