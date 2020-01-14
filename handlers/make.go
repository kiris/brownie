package handlers

import (
	"github.com/pkg/errors"

	"github.com/kiris/brownie/components"
	"github.com/kiris/brownie/lib/slack"
	"github.com/kiris/brownie/models"
)

type MakeCommandHandler struct {
	Renderer  components.Renderer
	Workspace *models.Workspace
}

func (h *MakeCommandHandler) ExecCommand(req *slack.CommandRequest) error {
	if len(req.CommandArgs) == 0 {
		repositories, err := h.Workspace.Repositories()
		if err != nil {
			return errors.Wrap(err, "failed to exec make command")
		}
		component := components.NewMakeComponent(req.Event.Channel)
		component.AppendSelectRepositoryAttachment(repositories)

		return h.Renderer.Render(component)
	} else {
		//// exec batch mode.
		//repoName := req.CommandArgs[0]
		//targets := req.CommandArgs[1:]
		//
		//repository := h.Workspace.Repository(repoName)
		//if repository == nil {
		//	return errors.Errorf("failed to exec make command. repository not found: name = %s", repoName)
		//}
		//
		//result := repository.RunMake(targets)
		//return h.sendResultMessages(req, result)
		return nil
	}
}

//func (h *MakeCommandHandler) sendResultMessages(req *CommandRequest, result *model.RunMakeResult) error {
//	ts, err := h.sendResultMessage(req, result)
//	if err != nil {
//		return err
//	}
//	return h.sendResultDetailMessage(req, ts, result)
//}
//func (h *MakeCommandHandler) sendResultMessage(req *CommandRequest, result *model.RunMakeResult) (string, error) {
//	user, err := h.Client.GetUserInfo(req.Event.User)
//	if err != nil {
//		return "", err
//	}
//
//	attachment := slack.Attachment{
//		Title: h.title(result),
//		Color: h.color(result),
//		Fields: []slack.AttachmentField{
//			{
//				Title: "project",
//				Value: result.Repository.Name,
//			},
//			{
//				Title: "branch",
//				Value: result.Branch,
//			},
//			{
//				Title: "targets",
//				Value: h.targets(result),
//			},
//		},
//		Footer: fmt.Sprintf("Executed by %s", user.Name),
//		FooterIcon: user.Profile.Image32,
//	}
//	options := slack.MsgOptionAttachments(attachment)
//
//	_, ts, err := h.Client.PostMessage(req.Event.Channel, options)
//	if err != nil {
//		return "", err
//	}
//
//	return ts, nil
//}
//
//func (h *MakeCommandHandler) title(result *model.RunMakeResult) string {
//	if result.Success {
//		return ":tada: make command SUCCESS!!"
//	} else {
//		return ":rain_cloud: make command FAILED..."
//	}
//}
//
//func (h *MakeCommandHandler) color(result *model.RunMakeResult) string {
//	if result.Success {
//		return "good"
//	} else {
//		return "danger"
//	}
//}
//
//func (h *MakeCommandHandler) targets(result *model.RunMakeResult) string {
//	if len(result.Targets) == 0 {
//		return "(default)"
//	} else {
//		return strings.Join(result.Targets, " ")
//	}
//}
//
//func (h *MakeCommandHandler) sendResultDetailMessage(req *CommandRequest, timestamp string, result *model.RunMakeResult) error {
//	ts := slack.MsgOptionTS(timestamp)
//	attachment := slack.MsgOptionAttachments(
//		slack.Attachment{
//			Title: ":memo: more details",
//			Color: "none",
//			Fields: []slack.AttachmentField{
//				{
//					Title: "exec command",
//					Value: result.Exec,
//				},
//				{
//					Title: "output",
//					Value: result.Output,
//				},
//			},
//		},
//	)
//
//	_, _, err := h.Client.PostMessage(req.Event.Channel, ts, attachment)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//
