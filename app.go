package brownie

import (
	"github.com/kiris/brownie/model"
	"github.com/nlopes/slack"

	"github.com/kiris/brownie/command"
	"github.com/kiris/brownie/interaction"
)

type App struct {
	client            *slack.Client
	commandListener   *command.Listener
	interactionServer *interaction.Server
	workspace         *model.Workspace
}

func CreateApp(slackToken, verificationToken, workSpaceDir string) *App {
	client := slack.New(slackToken)

	return &App{
		client           : client,
		commandListener  : command.CreateListener(client),
		interactionServer: interaction.CreateServer(verificationToken),
		workspace        : model.NewWorkspace(workSpaceDir),
	}
}

func (app *App) Run() error {
	app.commandListener.Handle("make", &command.MakeHandler{
		Client   : app.client,
		Workspace: app.workspace,
	})

	// TODO error code
	go app.commandListener.ListenAndResponse()
	//if err := app.commandListener.ListenAndResponse(); err != nil {
	//	return err
	//}

	app.interactionServer.Handle("cancel", &interaction.CancelHandler {
	})
	if err := app.interactionServer.ListenAndServ("8080"); err != nil {
		return err
	}


	return nil
}
