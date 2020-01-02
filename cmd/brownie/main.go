package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"

	"github.com/kiris/brownie/pkg/make"
)

type Env struct {
	SlackToken string `envconfig:"SLACK_TOKEN" required:"true"`
}

func init() {
	//// JSONフォーマット
	//log.SetFormatter(&log.JSONFormatter{})

	// 標準エラー出力でなく標準出力とする
	log.SetOutput(os.Stdout)

	// Warningレベル以上を出力
	log.SetLevel(log.DebugLevel)
}

func main() {
	var env Env
	if err := envconfig.Process("", &env); err != nil {
		log.WithFields(log.Fields{
			"msg": err,
		}).Error("Failed to process env.")
		os.Exit(1)
	}

	// Listening slack event and response
	log.Info("Start slack event listening")
	client := slack.New(env.SlackToken)

	_, err := client.GetChannels(false)
	if err != nil {
		log.WithFields(log.Fields{
			"msg": err,
		}).Error("Failed to slack api call.")
		os.Exit(1)
	}
	//for _, channel := range channels {
	//	log.WithFields(log.Fields{
	//		"id": channel.ID,
	//		"name": channel.Name,
	//	}).Info("")
	//
	//}


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
		os.Exit(1)
	}
}