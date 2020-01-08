package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/kiris/brownie"
)

func init() {
	//// JSONフォーマット
	//log.SetFormatter(&log.JSONFormatter{})

	// 標準エラー出力でなく標準出力とする
	log.SetOutput(os.Stdout)

	// Warningレベル以上を出力
	log.SetLevel(log.DebugLevel)
}

func main() {
	app, err := brownie.CreateAppFromEnvironmentVariables()
	if err != nil {
		log.WithError(err).Error("Failed to load env.")
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		log.WithError(err).Error("Failed to start listen and response.")
		os.Exit(1)
	}
}