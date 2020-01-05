package interaction

import (
	"encoding/json"
	"net/http"

	"github.com/nlopes/slack"
)
func responseMessage(w http.ResponseWriter, original slack.Message, title string, init bool) error {
	if init {
		original.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
	}

	if title != "" {
		original.Attachments[0].Fields = []slack.AttachmentField {
			{
				Title: title,
			},
		}
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(&original)
}
