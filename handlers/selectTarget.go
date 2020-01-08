package handlers

import (
	"github.com/kiris/brownie/components"
	"github.com/kiris/brownie/models"
	"github.com/nlopes/slack"
	"net/http"
)

type SelectTargetHandler struct {
	Workspace *models.Workspace
}

func (h *SelectTargetHandler) ServInteraction(w http.ResponseWriter, callback *slack.InteractionCallback) error {
	component := components.NewMakeComponentFromInteraction(callback, true)
	component.AppendConfirmExecAttachment()

	renderer := components.InteractionRenderer{
		Writer:   w,
		Callback: callback,
	}
	return renderer.Render(component)
}

