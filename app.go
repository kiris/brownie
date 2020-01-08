package brownie

import (
	"github.com/kiris/brownie/components"
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

func NewApp(slackToken, verificationToken, workSpaceDir string) *App {
	client := slack.New(slackToken)

	app := &App{
		client           : client,
		commandListener  : command.CreateListener(client),
		interactionServer: interaction.CreateServer(verificationToken, ":8081"),
		workspace        : model.NewWorkspace(workSpaceDir),
	}
	app.registerHandlers()

	return app
}

func (app *App) Run() error {
	errChan := make(chan error, 1)
	go func() {
		err := app.commandListener.ListenAndResponse()
		if err != nil {
			errChan <- errors.Wrap(err, "failed command listener start.")
		}
		errChan <- nil
	}()
	if err := app.interactionServer.ListenAndServ(); err != nil {
		return errors.Wrap(err, "failed interaction server start.")
	}

	return <-errChan
}


func (app *App) registerHandlers() {
	renderer := components.ApiRenderer{
		Client: app.client,
	}

	// commands
	app.commandListener.Handle("make", &command.MakeHandler{
		Client   : app.client,
		Renderer : renderer,
		Workspace: app.workspace,
	})

	// interactions
	app.interactionServer.Handle(components.ActionSelectRepository, &interaction.SelectRepositoryHandler {
		Workspace: app.workspace,
	})
	app.interactionServer.Handle(components.ActionSelectBranch, &interaction.SelectBranchHandler {
		Workspace: app.workspace,
	})
	app.interactionServer.Handle(components.ActionSelectTarget, &interaction.SelectTargetHandler {
		Workspace: app.workspace,
	})
	app.interactionServer.Handle(components.ActionExecMake, &interaction.ExecMakeHandler{
		Client:    app.client,
		Workspace: app.workspace,
	})
	app.interactionServer.Handle(components.ActionCancel, &interaction.CancelHandler {})
}
