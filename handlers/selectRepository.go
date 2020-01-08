package handlers

import (
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"net/http"

	"github.com/kiris/brownie/components"
	"github.com/kiris/brownie/models"
)

type SelectRepositoryHandler struct {
	Workspace *models.Workspace
}

func (h *SelectRepositoryHandler) ServInteraction(w http.ResponseWriter, callback *slack.InteractionCallback) error {
	component := components.NewMakeComponentFromInteraction(callback, true)

	repoName := component.SelectedRepository()
	repository := h.Workspace.GetRepository(repoName)
	if repository == nil {
		return errors.Errorf("failed to exec make command. repository not found: name = %s", repoName)
	}
	component.AppendSelectBranchAttachment()

	renderer := components.InteractionRenderer{
		Writer:   w,
		Callback: callback,
	}
	return renderer.Render(component)
}

