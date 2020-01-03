package brownie

import (
	"fmt"
)

type App struct {
	slack     *SlackListener
	workspace *Workspace
}

func CreateAppFromEnvironmentVariables() (*App, error) {
	env, err := LoadEnv()
	if err != nil {
		return nil, err
	}

	return &App{
		slack:     NewSlackListener(env.SlackToken),
		workspace: NewWorkspace(env.WorkspaceDir),
	}, nil
}


func (app *App) StartListenAndResponse() error {
	return app.slack.ListenAndResponse(app)
}

func (app *App) ExecMake(projectName string, targets []string) (*MakeResult, error) {
	project := app.workspace.GetProject(projectName)

	if project == nil {
		return nil, fmt.Errorf("project not found. name = %s", projectName)
	}

	return project.ExecMake(targets), nil
}