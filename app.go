package brownie

import (
	"github.com/nlopes/slack"
	"github.com/pkg/errors"

	"github.com/kiris/brownie/command"
	"github.com/kiris/brownie/interaction"
	"github.com/kiris/brownie/model"
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
		interactionServer: interaction.CreateServer(verificationToken, ":8081"),
		workspace        : model.NewWorkspace(workSpaceDir),
	}
}

func (app *App) Run() error {
	app.commandListener.Handle("make", &command.MakeHandler{
		Client   : app.client,
		Workspace: app.workspace,
	})

	errChan := make(chan error, 1)
	go func() {
		err := app.commandListener.ListenAndResponse()
		if err != nil {
			errChan <- errors.Wrap(err, "failed command listener start.")
		}
		errChan <- nil
	}()

	app.interactionServer.Handle(interaction.ActionSelectRepository, &interaction.SelectRepositoryHandler {
		Workspace: app.workspace,
	})
	app.interactionServer.Handle(interaction.ActionSelectBranch, &interaction.SelectBranchHandler {
		Workspace: app.workspace,
	})
	app.interactionServer.Handle(interaction.ActionSelectTarget, &interaction.SelectTargetHandler {
		Workspace: app.workspace,
	})
	app.interactionServer.Handle(interaction.ActionExecMake, &interaction.ExecMakeHandler{
		Client:    app.client,
		Workspace: app.workspace,
	})
	app.interactionServer.Handle(interaction.ActionCancel, &interaction.CancelHandler {})
	if err := app.interactionServer.ListenAndServ(); err != nil {
		return errors.Wrap(err, "failed interaction server start.")
	}

	return <-errChan
}
