package interaction

import (
	"encoding/json"
	"github.com/kiris/brownie/model"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"net/http"
)

type SelectBranchHandler struct {
	Workspace *model.Workspace
}

func (h *SelectBranchHandler) ServInteraction(w http.ResponseWriter, callback *slack.InteractionCallback) error {
	component := NewMakeSettingsComponentFromCallback(callback, true)

	repoName := component.GetSelectedRepository()
	repository := h.Workspace.GetRepository(repoName)
	if repository == nil {
		return errors.Errorf("failed to exec make command. repository not found: name = %s", repoName)
	}
	targets, _ := repository.CollectMakeTargets()
	component.AppendSelectTargetAttachment(targets)

	original := callback.OriginalMessage
	original.ReplaceOriginal = true
	original.Attachments = component.Attachments
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(&original)
}
