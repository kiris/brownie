package interaction

import (
	"encoding/json"
	"github.com/kiris/brownie/model"
	"github.com/nlopes/slack"
	"net/http"
)

type SelectTargetHandler struct {
	Workspace *model.Workspace
}

func (h *SelectTargetHandler) ServInteraction(w http.ResponseWriter, callback *slack.InteractionCallback) error {
	component := NewMakeSettingsComponentFromCallback(callback, true)
	component.AppendExecMakeAttachment()

	original := callback.OriginalMessage
	original.ReplaceOriginal = true
	original.Attachments = component.Attachments
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(&original)
}

