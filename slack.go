package brownie

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

type SlackListener struct {
	client  *slack.Client
	rtm     *slack.RTM
	botId   string
	botName string
}

type IResponse interface {
	Send(rtm *slack.RTM) error
}

func NewSlackListener(token string) *SlackListener {
	client := slack.New(
		token,
	)

	return &SlackListener{
		client: client,
		rtm: client.NewRTM(),
	}
}

func (s *SlackListener) ListenAndResponse(b *App) error {
	go s.rtm.ManageConnection()

	for msg := range s.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectingEvent:
			log.Info("Slack RTM connecting.")

		case *slack.ConnectedEvent:
			log.WithFields(log.Fields{
				"id":   ev.Info.User.ID,
				"name": ev.Info.User.Name,
			}).Info("Slack RTM connected.")

			s.botId = ev.Info.User.ID
			s.botName = ev.Info.User.Name

		case *slack.MessageEvent:
			s.handleMessageEvent(ev, b)

		case *slack.InvalidAuthEvent:
			return errors.New("invalid credentials")

		default:
			// Ignore other events.
			// fmt.Printf("other: %v\n", ev)
		}
	}

	return nil
}

func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent, app *App) {
	tokens := s.messageTokens(ev)

	// Only response mention to bot. Ignore else.
	if !s.mentionToMe(tokens) {
		return
	}

	if s.messageFromBot(ev) {
		log.Debug("ignore bot message.")
		return
	}

	// @brownie [command] arg1 arg2 ...
	cmd := s.getCommand(tokens)
	cmdArgs := s.getCommandArgs(tokens)

	response, err := s.runCommand(ev, app, cmd, cmdArgs)
	if err != nil {
		// TODO logging or response.
		return
	}

	if err := response.Send(s.rtm); err != nil {
		// TODO logging or response.
		return
	}
}

func (s *SlackListener) getCommand(tokens []string) string {
	if len(tokens) < 2 {
		return "help"
	}

	return tokens[1]
}

func (s *SlackListener) getCommandArgs(tokens []string) []string {
	return tokens[2:]
}

func (s *SlackListener) runCommand(ev *slack.MessageEvent, app *App, cmd string, cmdArgs []string) (IResponse, error) {
	switch cmd {
	case "make":
		if len(cmdArgs) == 0 {
			// exec interactive mode.
			return nil, errors.New("interactive mode is not implemented")
		} else {
			// exec batch mode.
			project := cmdArgs[0]
			targets := cmdArgs[1:]

			result, err := app.ExecMake(project, targets)
			if  err != nil {
				return nil, err
			}

			return &ExecMakeResponse{
				event: ev,
				result: result,
			}, nil
		}

	default:
		// unknown command.
		return nil, fmt.Errorf("%s is unknown command", cmd)
	}
}

func (s *SlackListener) messageTokens(ev *slack.MessageEvent) []string {
	return strings.Split(strings.TrimSpace(ev.Msg.Text), " ")
}

func (s *SlackListener) mentionToMe(tokens []string) bool {
	return tokens[0] == fmt.Sprintf("<@%s>", s.botId)
}

func (s *SlackListener) messageFromBot(ev *slack.MessageEvent) bool {
	return ev.Msg.SubType == "bot_message"
}


