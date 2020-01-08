package interaction

import (
	"encoding/json"
	"net/http"

	"github.com/nlopes/slack"
)

type CancelHandler struct {
}

func (h *CancelHandler) ServInteraction(w http.ResponseWriter, callback *slack.InteractionCallback) error {
	component := NewMakeSettingsComponentFromCallback(callback, false)
	component.Cancel(callback.User)

	original := callback.OriginalMessage
	original.ReplaceOriginal = true
	original.Attachments = component.Attachments
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(&original)
}