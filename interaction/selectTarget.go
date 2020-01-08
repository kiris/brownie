package interaction

import (
	"github.com/kiris/brownie/components"
	"github.com/kiris/brownie/model"
	"github.com/nlopes/slack"
	"net/http"
)

type SelectTargetHandler struct {
	Workspace *model.Workspace
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

