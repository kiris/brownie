package command

import (
	"github.com/kiris/brownie/components"
	"github.com/pkg/errors"

	"github.com/kiris/brownie/model"
)

type MakeHandler struct {
	Renderer  components.Renderer
	Workspace *model.Workspace
}

func (h *MakeHandler) ExecCommand(req *Request) error {
	if len(req.CommandArgs) == 0 {
		repositories, err := h.Workspace.GetRepositories()
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
		//repository := h.Workspace.GetRepository(repoName)
		//if repository == nil {
		//	return errors.Errorf("failed to exec make command. repository not found: name = %s", repoName)
		//}
		//
		//result := repository.ExecMake(targets)
		//return h.sendResultMessages(req, result)
		return nil
	}
}

//func (h *MakeHandler) sendResultMessages(req *Request, result *model.ExecMakeResult) error {
//	ts, err := h.sendResultMessage(req, result)
//	if err != nil {
//		return err
//	}
//	return h.sendResultDetailMessage(req, ts, result)
//}
//func (h *MakeHandler) sendResultMessage(req *Request, result *model.ExecMakeResult) (string, error) {
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
//func (h *MakeHandler) title(result *model.ExecMakeResult) string {
//	if result.Success {
//		return ":tada: make command SUCCESS!!"
//	} else {
//		return ":rain_cloud: make command FAILED..."
//	}
//}
//
//func (h *MakeHandler) color(result *model.ExecMakeResult) string {
//	if result.Success {
//		return "good"
//	} else {
//		return "danger"
//	}
//}
//
//func (h *MakeHandler) targets(result *model.ExecMakeResult) string {
//	if len(result.Targets) == 0 {
//		return "(default)"
//	} else {
//		return strings.Join(result.Targets, " ")
//	}
//}
//
//func (h *MakeHandler) sendResultDetailMessage(req *Request, timestamp string, result *model.ExecMakeResult) error {
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
