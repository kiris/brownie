package handlers

import (
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/kiris/brownie/components"
	"github.com/kiris/brownie/models"
)

type ExecMakeHandler struct {
	ApiRenderer *components.ApiRenderer
	Workspace   *models.Workspace
}

func (h *ExecMakeHandler) ServInteraction(w http.ResponseWriter, callback *slack.InteractionCallback) error {
	component := components.NewMakeComponentFromInteraction(callback, false)

	selectedRepository := component.SelectedRepository()
	selectedBranch := component.SelectedBranch()
	selectedTarget := component.SelectedTarget()
	component.InProgress(callback.User)

	renderer := components.InteractionRenderer{
		Writer:   w,
		Callback: callback,
	}

	if err := renderer.Render(component); err != nil {
		return err
	}

	go func() {
		result, err := h.execMake(selectedRepository, selectedBranch, selectedTarget)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"repository": selectedRepository,
				"branch":     selectedBranch,
				"target":     selectedTarget,
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
	if err := h.ApiRenderer.Render(component); err != err {
		return errors.Wrap(err, "failed to async response.")
	}

	if err := h.ApiRenderer.Render(outputComponent); err != err {
		return errors.Wrap(err, "failed to async response.")
	}

	return nil
}

func (h *ExecMakeHandler) execMake(repoName string, branchName string, target string) (*models.RunMakeResult, error) {
	repository := h.Workspace.Repository(repoName)
	if repository == nil {
		return nil, errors.Errorf("repository not found: name = %s", repoName)
	}

	return repository.RunMake([]string {target}), nil
}

