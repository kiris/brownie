package interaction

import (
	"fmt"
	"github.com/nlopes/slack"
	"net/http"
)

type CancelHandler struct {
}

func (h *CancelHandler) ServInteraction(w http.ResponseWriter, message slack.InteractionCallback) error {
	title := fmt.Sprintf(":x: %s canceled the command", message.User.Name)
	return responseMessage(w, message.OriginalMessage, title, true)
}