package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"

	"github.com/kiris/brownie/pkg/make"
)

type Env struct {
	SlackToken string `envconfig:"SLACK_TOKEN" required:"true"`
}

func main() {
	var env Env
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	// Listening slack event and response
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(env.SlackToken)

	cannels, err := client.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	for _, channel := range cannels {
		fmt.Printf("ID: %s, Name: %s\n", channel.ID, channel.Name)
	}

	out, err := make.ExecMake("workspace/kiribot", "master", "usage")

	fmt.Print(string(out))
	fmt.Print(err)
}