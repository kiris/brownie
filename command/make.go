package command

import (
	"fmt"
	"strings"

	"github.com/kiris/brownie/model"
	"github.com/nlopes/slack"
)

type MakeHandler struct {
	Client    *slack.Client
	Workspace *model.Workspace
}

func (h *MakeHandler) ExecCommand(req *Request) error {
	if len(req.CommandArgs) == 0 {
		// exec interactive mode.
		return h.postConfigMakeMessage(req)
	} else {
		// exec batch mode.
		projectName := req.CommandArgs[0]
		targets := req.CommandArgs[1:]

		result, err := h.execMake(projectName, targets)
		if err != nil {
			// TODO
			return err
		}

		return h.postResultMessages(req, result)
	}
}


func (h *MakeHandler) execMake(projectName string, targets []string) (*model.ExecMakeResult, error) {
	project := h.Workspace.GetProject(projectName)
	if project == nil {
		return nil, fmt.Errorf("project not found. name = %s", projectName)
	}

	result := project.ExecMake(targets)
	return result, nil
}



func (h *MakeHandler) postResultMessages(cmd *Request, result *model.ExecMakeResult) error {
	if ts, err := h.sendResultMessage(cmd, result); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	} else {
		if _, err := h.sentResultDetailMessage(cmd, ts, result); err != nil {
			return fmt.Errorf("failed to post message: %s", err)
		}
	}

	return nil
}

func (h *MakeHandler) sendResultMessage(cmd *Request, result *model.ExecMakeResult) (string, error) {
	user, err := h.Client.GetUserInfo(cmd.Event.User)
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %s", cmd.Event.User)
	}

	attachment := slack.Attachment{
		Title: h.title(result),
		Color: h.color(result),
		Fields: []slack.AttachmentField{
			{
				Title: "project",
				Value: result.Project.Name,
			},
			{
				Title: "branch",
				Value: result.Branch,
			},
			{
				Title: "targets",
				Value: h.targets(result),
			},
		},
		Footer: fmt.Sprintf("Executed by %s", user.Name),
		FooterIcon: user.Profile.Image32,
	}

	return cmd.ResponseAttachmentsMessage(attachment)
}

func (h *MakeHandler) title(result *model.ExecMakeResult) string {
	if result.Success {
		return ":tada: make command SUCCESS!!"
	} else {
		return ":rain_cloud: make command FAILED..."
	}
}

func (h *MakeHandler) color(result *model.ExecMakeResult) string {
	if result.Success {
		return "good"
	} else {
		return "danger"
	}
}

func (h *MakeHandler) targets(result *model.ExecMakeResult) string {
	if len(result.Targets) == 0 {
		return "(default)"
	} else {
		return strings.Join(result.Targets, " ")
	}
}

func (h *MakeHandler) sentResultDetailMessage(cmd *Request, timestamp string, result *model.ExecMakeResult) (string, error) {
	ts := slack.MsgOptionTS(timestamp)
	attachment := slack.MsgOptionAttachments(
		slack.Attachment{
			Title: ":memo: more details",
			Color: "none",
			Fields: []slack.AttachmentField{
				{
					Title: "exec command",
					Value: result.Exec,
				},
				{
					Title: "output",
					Value: result.Output,
				},
			},
		},
	)

	return cmd.ResponseMessage(ts, attachment)
}


func (h *MakeHandler) postConfigMakeMessage(cmd *Request) error {
	attachment := slack.Attachment{
		Title     : "select branch",
		CallbackID: "make",
		Actions   : []slack.AttachmentAction{
			{
				Name   : "branch",
				Type   : "select",
				Options: []slack.AttachmentActionOption {
					{
						Text : "master",
						Value: "master",
					},
					{
						Text : "staging",
						Value: "staging",
					},

				},
			},
			{
				Name : "cancel",
				Text : "Cancel",
				Type : "button",
				Style: "danger",
			},
		},
	}

	if _, err := cmd.ResponseAttachmentsMessage(attachment); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}

