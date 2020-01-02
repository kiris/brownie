package brownie

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"

	"github.com/kiris/brownie/pkg/make"
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

func (server *Server) Start() error {
	server.rtm = server.client.NewRTM().NewRTM()
	go server.rtm.ManageConnection()

	for msg := range server.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			log.WithFields(log.Fields{
				"id":   ev.Info.User.ID,
				"name": ev.Info.User.Name,
			}).Info("success connection.")
			server.botId = ev.Info.User.ID
			server.botName = ev.Info.User.Name

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)
			server.handleMessageEvent(ev)

		case *slack.InvalidAuthEvent:
			return errors.New("invalid credentials")

		default:
			// Ignore other events..
			fmt.Printf("other: %v\n", ev)
		}
	}

	return nil
}

func (server *Server) handleMessageEvent(ev *slack.MessageEvent) {
	// Only response mention to bot. Ignore else.
	msg := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")


	log.WithFields(log.Fields{
		"msg": msg,
	}).Info("mention message")
	if msg[0] != fmt.Sprintf("<@%s>", server.botId) {
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
		currentDir, _ := os.Getwd()
		cmd := make.Make {
			Dir: currentDir + "/workspace/kiribot",
			Branch: "master",
			Targets: []string { "usage" },
			DryRun: false,
		}
		out, err := cmd.Exec()

		log.Info(string(out))

		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Failed to exec make.")
		}

		//args := msg[2:]
		//if err := s.ResponseDeploy(ev); err != nil {
		//	return fmt.Errorf("failed to post message: %s", err)
		//}
	}



}