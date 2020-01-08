package brownie

import (
	"github.com/kiris/brownie/components"
	"github.com/kiris/brownie/handlers"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"

	libSlack "github.com/kiris/brownie/lib/slack"
	"github.com/kiris/brownie/models"
)

type App struct {
	client            *slack.Client
	commandListener   *libSlack.CommandListener
	interactionServer *libSlack.InteractionServer
	workspace         *models.Workspace
}

func NewApp(slackToken, verificationToken, workSpaceDir string) *App {
	client := slack.New(slackToken)

	app := &App{
		client           : client,
		commandListener  : libSlack.NewCommandListener(client),
		interactionServer: libSlack.NewInteractionServer(verificationToken, ":8081"),
		workspace        : models.NewWorkspace(workSpaceDir),
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
	app.commandListener.Handle("make", &handlers.MakeHandler{
		Renderer : &renderer,
		Workspace: app.workspace,
	})

	// interactions
	app.interactionServer.Handle(components.ActionSelectRepository, &handlers.SelectRepositoryHandler {
		Workspace: app.workspace,
	})
	app.interactionServer.Handle(components.ActionSelectBranch, &handlers.SelectBranchHandler {
		Workspace: app.workspace,
	})
	app.interactionServer.Handle(components.ActionSelectTarget, &handlers.SelectTargetHandler {
		Workspace: app.workspace,
	})
	app.interactionServer.Handle(components.ActionExecMake, &handlers.ExecMakeHandler{
		Client:    app.client,
		Workspace: app.workspace,
	})
	app.interactionServer.Handle(components.ActionCancel, &handlers.CancelHandler {})
}
