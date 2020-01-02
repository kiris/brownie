package brownie

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	client  *slack.Client
	rtm     *slack.RTM
	botId   string
	botName string
}


func NewServer(token string) *Server {
	client := slack.New(
		token,
	)

	return &Server {
		client: client,
	}
}

func (s *Server) Start(b *Brownie) error {
	s.rtm = s.client.NewRTM().NewRTM()
	go s.rtm.ManageConnection()

	for msg := range s.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			log.WithFields(log.Fields{
				"id":   ev.Info.User.ID,
				"name": ev.Info.User.Name,
			}).Info("success connection.")
			s.botId = ev.Info.User.ID
			s.botName = ev.Info.User.Name

		case *slack.MessageEvent:
			s.handleMessageEvent(ev, b)

		case *slack.InvalidAuthEvent:
			return errors.New("invalid credentials")

		default:
			// Ignore other events..
			// fmt.Printf("other: %v\n", ev)
		}
	}

	return nil
}

func (s *Server) handleMessageEvent(ev *slack.MessageEvent, b *Brownie) {
	// Only response mention to bot. Ignore else.
	msg := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")

	if msg[0] != fmt.Sprintf("<@%s>", s.botId) || len(msg) < 2 {
		return
	}

	if ev.Msg.SubType == "bot_message" {
		log.Debug("slack/ignore bot message")
		return
	}

	log.WithFields(log.Fields{
		"msg": msg,
	}).Info("mention message")

	switch msg[1] {
	case "make":
		if (len(msg) == 2) {
			// exec interactive mode.
		} else {
			// exec batch mode.
			project := msg[2]
			targets := msg[3:]
			b.ExecMake(project, targets)
		}
	default:
		log.WithFields(log.Fields{
			"msg": msg[1],
		}).Info("unknown message")

	}
}

type RunMakeArgs struct {
	project string
	targets []string
}