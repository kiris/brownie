package interaction

import (
	"fmt"
	"github.com/kiris/brownie/model"
	"github.com/nlopes/slack"
	"strings"
)


type MakeSettingsComponent struct {
	channel string
	ts string
	Attachments []slack.Attachment
}

func NewMakeSettingsComponentFromCallback(callback *slack.InteractionCallback, changeSetting bool) MakeSettingsComponent {
	component := MakeSettingsComponent{
		channel:     callback.Channel.ID,
		ts:          callback.MessageTs,
		Attachments: callback.OriginalMessage.Attachments,
	}

	if changeSetting {
		component.changeSetting(callback.ActionCallback.AttachmentActions[0])
	}

	return component
}

func (c *MakeSettingsComponent) GetSelectedRepository() string {
	return c.Attachments[0].Actions[0].SelectedOptions[0].Value
}

func (c *MakeSettingsComponent) GetSelectedBranch() string {
	return c.Attachments[1].Actions[0].SelectedOptions[0].Value
}

func (c *MakeSettingsComponent) GetSelectedTarget() string {
	return c.Attachments[2].Actions[0].SelectedOptions[0].Value
}

func (c *MakeSettingsComponent) AppendSelectRepositoryAttachment(repositories []*model.Repository) {
	c.Attachments = append(c.Attachments, c.newSelectRepositoryAttachment(repositories))
}

func (c *MakeSettingsComponent) AppendSelectBranchAttachment() {
	c.Attachments = append(c.Attachments, c.newSelectBranchAttachment(nil))
}

func (c *MakeSettingsComponent) AppendSelectTargetAttachment(targets []string) {
	c.Attachments = append(c.Attachments, c.newSelectTargetAttachment(targets))
}

func (c *MakeSettingsComponent) AppendExecMakeAttachment() {
	c.Attachments = append(c.Attachments, c.newExecMakeAttachment())
}


func (c *MakeSettingsComponent) InProgress(user slack.User) {
	repository := c.GetSelectedRepository()
	branch := c.GetSelectedBranch()
	target := c.GetSelectedTarget()

	c.Attachments = []slack.Attachment {
		{
			Color: "warning",
			Title: fmt.Sprintf(":hammer_and_wrench: The Make command is running ..."),
			Text: fmt.Sprintf("REPOSITORY: %s\nBRANCH: %s\nTARGET: %s", repository, branch, target),
			Footer: fmt.Sprintf("executed by %s", user.Name),
			FooterIcon: user.Profile.Image32,
		},
	}
}

func (c *MakeSettingsComponent) Done(result *model.ExecMakeResult) []slack.Attachment {
	if result.Success {
		c.Attachments[0].Title = ":tada: make command SUCCESS!!"
		c.Attachments[0].Color = "good"
	} else {
		c.Attachments[0].Title = ":rain_cloud: make command FAILED..."
		c.Attachments[0].Color = "error"
	}

	detailAttachments := []slack.Attachment{
		{
			Title: ":memo: details",
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
	}

	return detailAttachments
}


func (c *MakeSettingsComponent) Cancel(user slack.User) {
	c.Attachments = []slack.Attachment {
		{
			Color: "danger",
			Fields: []slack.AttachmentField{
				{
					Title: fmt.Sprintf(":x: This request has been canceled"),
				},
			},
			Footer: fmt.Sprintf("canceled by %s", user.Name),
			FooterIcon: user.Profile.Image32,
		},
	}
}



func (c *MakeSettingsComponent) changeSetting(selectedAction *slack.AttachmentAction) {

	selectedOptionValue := selectedAction.SelectedOptions[0].Value

	for i := range c.Attachments {
		attachment := &c.Attachments[i]

		if attachment.Actions[0].Name == selectedAction.Name {
			for _, option := range attachment.Actions[0].Options {
				if option.Value == selectedOptionValue {
					// selected
					attachment.Color = "good"
					attachment.Actions = attachment.Actions[0:1] // remove cancel button
					attachment.Actions[0].SelectedOptions = []slack.AttachmentActionOption{option}
					c.Attachments = c.Attachments[0:i+1]
					return
				}
			}
		}
	}
}

func (c *MakeSettingsComponent) newSelectRepositoryAttachment(repositories []*model.Repository) slack.Attachment {
	options := make([]slack.AttachmentActionOption, len(repositories))
	for i, r := range repositories {
		options[i] = slack.AttachmentActionOption{
			Text:  r.Name,
			Value: r.Name,
		}
	}

	attachment := slack.Attachment{
		Title:      "select repository",
		CallbackID: "repository",
		Actions:    []slack.AttachmentAction{
			{
				Name:    ActionSelectRepository,
				Text:    "Choose a repository ...",
				Type:    "select",
				Options: options,
			},
			c.newCancelAction(),
		},
	}

	return attachment
}

func (c *MakeSettingsComponent) newSelectBranchAttachment(_ []string) slack.Attachment {
	options := []slack.AttachmentActionOption{
		{
			Text:  "master",
			Value: "master",
		},
		{
			Text:  "staging",
			Value: "staging",
		},
	}

	attachment := slack.Attachment{
		Title:      "select branch",
		CallbackID: "branch",
		Actions:    []slack.AttachmentAction{
			{
				Name:    ActionSelectBranch,
				Text:    "Choose a branch ...",
				Type:    "select",
				Options: options,
			},
			c.newCancelAction(),
		},
	}

	return attachment
}

func (c *MakeSettingsComponent) newSelectTargetAttachment(targets []string) slack.Attachment {
	options := make([]slack.AttachmentActionOption, len(targets))
	for i, target := range targets {
		options[i] = slack.AttachmentActionOption{
			Text:  target,
			Value: target,
		}
	}

	attachment := slack.Attachment{
		Title:      "select target",
		CallbackID: "target",
		Actions:    []slack.AttachmentAction{
			{
				Name:    ActionSelectTarget,
				Text:    "Choose a target ...",
				Type:    "select",
				Options: options,
			},
			c.newCancelAction(),
		},
	}

	return attachment
}

func (c *MakeSettingsComponent) newExecMakeAttachment() slack.Attachment {
	attachment := slack.Attachment{
		Title:      "exec make?",
		CallbackID: "exec",
		Actions:    []slack.AttachmentAction{
			{
				Name:  ActionExecMake,
				Text:  "Exec Make",
				Type:  "button",
				Style: "primary",
			},
			c.newCancelAction(),
		},
	}

	return attachment
}

func (c *MakeSettingsComponent) newCancelAction() slack.AttachmentAction {
	return slack.AttachmentAction{
		Name:  ActionCancel,
		Text:  "Cancel",
		Type:  "button",
		Style: "danger",
	}
}




func (c *MakeSettingsComponent) targets(result *model.ExecMakeResult) string {
	if len(result.Targets) == 0 {
		return "(default)"
	} else {
		return strings.Join(result.Targets, " ")
	}
}
