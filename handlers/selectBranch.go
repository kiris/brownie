package handlers

import (
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"net/http"

	"github.com/kiris/brownie/components"
	"github.com/kiris/brownie/models"
)

type SelectBranchHandler struct {
	Workspace *models.Workspace
}

func (h *SelectBranchHandler) ServInteraction(w http.ResponseWriter, callback *slack.InteractionCallback) error {
	component := components.NewMakeComponentFromInteraction(callback, true)

	selectedRepository := component.SelectedRepository()
	repository := h.Workspace.Repository(selectedRepository)
	if repository == nil {
		return errors.Errorf("failed to exec make command. repository not found: name = %s", selectedRepository)
	}
	targets, _ := repository.Targets()
	component.AppendSelectTargetAttachment(targets)

	renderer := components.InteractionRenderer{
		Writer:   w,
		Callback: callback,
	}
	return renderer.Render(component)
}
