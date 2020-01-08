package interaction

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/kiris/brownie/model"
)

type ExecMakeHandler struct {
	Client    *slack.Client
	Workspace *model.Workspace
}

func (h *ExecMakeHandler) ServInteraction(w http.ResponseWriter, callback *slack.InteractionCallback) error {
	component := NewMakeSettingsComponentFromCallback(callback, false)

	repoName := component.GetSelectedRepository()
	branchName := component.GetSelectedBranch()
	target := component.GetSelectedTarget()

	component.InProgress(callback.User)
	original := callback.OriginalMessage
	original.ReplaceOriginal = true
	original.Attachments = component.Attachments
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&original); err != nil {
		return err
	}

	go func() {
		result, err := h.execMake(repoName, branchName, target)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"repository": repoName,
				"branch": branchName,
				"target": target,
			}).Error("failed to exec make command.")
		}
		detailAttachments := component.Done(result)

		options := slack.MsgOptionAttachments(component.Attachments ...)

		if _, _, _, err := h.Client.UpdateMessage(component.channel, component.ts, options); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"channel": component.channel,
				"ts": component.ts,
				"options": options,
			}).Error( "failed to update message.")
		}

		tsOption := slack.MsgOptionTS(component.ts)
		detailAttachmentsOption := slack.MsgOptionAttachments(detailAttachments ...)
		if _, _, _, err := h.Client.SendMessage(component.channel, tsOption, detailAttachmentsOption); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"channel": component.channel,
				"ts": component.ts,
				"options": options,
			}).Error( "failed to send detail message.")
		}


	}()


	return nil
}

func (h *ExecMakeHandler) execMake(repoName string, branchName string, target string) (*model.ExecMakeResult, error) {
	repository := h.Workspace.GetRepository(repoName)
	if repository == nil {
		return nil, errors.Errorf("repository not found: name = %s", repoName)
	}

	return repository.ExecMake([]string {target}), nil
}

