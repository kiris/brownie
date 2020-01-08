package handlers

import (
	"regexp"


	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"

	libSlack "github.com/kiris/brownie/lib/slack"
	"github.com/kiris/brownie/models"
)

type CloneCommandHandler struct {
	Client *slack.Client
	Workspace *models.Workspace
}

var formattedUrlRegexp = regexp.MustCompile(`^<(.+)>$`)
var extractRepositoryNameRegexp = regexp.MustCompile(`^.+/([^/]+)/?>$`)
func (h *CloneCommandHandler) ExecCommand(req *libSlack.CommandRequest) error {
	if len(req.CommandArgs) != 1 {
		return nil
	}

	url := h.unFormatUrl(req.CommandArgs[0])
	_, err := h.Workspace.CreateRepository(url)

	attachment := slack.Attachment{}
	if err == nil {
		attachment.Title = ":tada: clone command SUCCESS!!"
		attachment.Color = "good"
	} else {
		log.WithError(err).WithField("url", url).Warn("failed to get command")
		attachment.Title = ":rain_cloud: clone command FAILED..."
		attachment.Color = "danger"
		attachment.Text = err.Error() // TODO error handling
	}

	attachmentsOption := slack.MsgOptionAttachments(attachment)
	if _, _, err = h.Client.PostMessage(req.Event.Channel, attachmentsOption); err != nil {
		return err
	}

	return err
}

func (h *CloneCommandHandler) unFormatUrl(url string) string {
	results := formattedUrlRegexp.FindStringSubmatch(url)

	if len(results) == 2 {
		return results[1]
	} else {
		return url
	}
}

