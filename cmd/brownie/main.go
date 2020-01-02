package main

import (
	"github.com/kelseyhightower/envconfig"
	brownie2 "github.com/kiris/brownie"
	log "github.com/sirupsen/logrus"
	"os"
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

	server := brownie2.NewServer(env.SlackToken)
	if err := server.Start(); err != nil {
		log.WithFields(log.Fields{
			"msg": err,
		}).Error("Failed to server start.")
		os.Exit(1)
	}
}