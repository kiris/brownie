package interaction

import (
	"github.com/kiris/brownie/components"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/kiris/brownie/model"
)

type ExecMakeHandler struct {
	Client    *slack.Client
	Workspace *model.Workspace
}

func (h *ExecMakeHandler) ServInteraction(w http.ResponseWriter, callback *slack.InteractionCallback) error {
	component := components.NewMakeComponentFromInteraction(callback, false)

	repoName := component.SelectedRepository()
	branchName := component.SelectedBranch()
	target := component.SelectedTarget()
	component.InProgress(callback.User)

	renderer := components.InteractionRenderer{
		Writer:   w,
		Callback: callback,
	}

	if err := renderer.Render(component); err != nil {
		return err
	}

	go func() {
		result, err := h.execMake(repoName, branchName, target)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"repository": repoName,
				"branch": branchName,
				"target": target,
			}).Error("failed to async exec make.")
		}
		outputComponent := component.Done(result)
		if err := h.asyncResponse(component, outputComponent); err != nil {
			log.WithError(err).Error("failed to async exec make.")
		}
	}()

	return nil
}

func (h *ExecMakeHandler) asyncResponse(component *components.MakeComponent, outputComponent *components.MakeOutputComponent) error {
	renderer := components.ApiRenderer{
		Client: h.Client,
	}
	if err := renderer.Render(component); err != err {
		return errors.Wrap(err, "failed to async response.")
	}

	if err := renderer.Render(outputComponent); err != err {
		return errors.Wrap(err, "failed to async response.")
	}

	return nil
}

func (h *ExecMakeHandler) execMake(repoName string, branchName string, target string) (*model.ExecMakeResult, error) {
	repository := h.Workspace.GetRepository(repoName)
	if repository == nil {
		return nil, errors.Errorf("repository not found: name = %s", repoName)
	}

	return repository.ExecMake([]string {target}), nil
}

