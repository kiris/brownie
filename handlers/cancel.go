package handlers

import (
	"net/http"

	"github.com/nlopes/slack"

	"github.com/kiris/brownie/components"
)

type CancelHandler struct {
}

func (h *CancelHandler) ServInteraction(w http.ResponseWriter, callback *slack.InteractionCallback) error {
	component := components.NewMakeComponentFromInteraction(callback, false)
	component.Cancel(callback.User)

	renderer := components.InteractionRenderer{
		Writer:   w,
		Callback: callback,
	}
	return renderer.Render(component)
}