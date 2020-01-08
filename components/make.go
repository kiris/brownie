package components

import (
	"fmt"
	"github.com/nlopes/slack"

	"github.com/kiris/brownie/models"
)

const (
	ActionSelectRepository = "selectRepository"
	ActionSelectBranch = "selectBranch"
	ActionSelectTarget = "selectTarget"
	ActionExecMake = "execMake"
	ActionCancel = "cancel"
)

const (
	selectRepository = 0
	selectBranch
	selectTarget
	confirmExec
	inProgress
	success
	failed
	cancel
)

type MakeComponent struct {
	Channel     string
	Ts          string
	attachments []slack.Attachment
}

func NewMakeComponent(channel string) *MakeComponent {
	return &MakeComponent{
		Channel: channel,
	}
}
func NewMakeComponentFromInteraction(callback *slack.InteractionCallback, changeSetting bool) *MakeComponent {
	component := &MakeComponent{
		Channel:     callback.Channel.ID,
		Ts:          callback.MessageTs,
		attachments: callback.OriginalMessage.Attachments,
	}

	if changeSetting {
		component.changeSetting(callback.ActionCallback.AttachmentActions[0])
	}

	return component
}

func (c *MakeComponent) Render() (string, string, string, []slack.Attachment) {
	return c.Channel, c.Ts, "", c.attachments
}

func (c *MakeComponent) SelectedRepository() string {
	return c.attachments[0].Actions[0].SelectedOptions[0].Value
}

func (c *MakeComponent) SelectedBranch() string {
	return c.attachments[1].Actions[0].SelectedOptions[0].Value
}

func (c *MakeComponent) SelectedTarget() string {
	return c.attachments[2].Actions[0].SelectedOptions[0].Value
}

func (c *MakeComponent) AppendSelectRepositoryAttachment(repositories []*models.Repository) {
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

	c.attachments = append(c.attachments, attachment)
}

func (c *MakeComponent) AppendSelectBranchAttachment() {
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

	c.attachments = append(c.attachments, attachment)
}

func (c *MakeComponent) AppendSelectTargetAttachment(targets []string) {
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

	c.attachments = append(c.attachments, attachment)
}

func (c *MakeComponent) AppendConfirmExecAttachment() {
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

	c.attachments = append(c.attachments, attachment)
}


func (c *MakeComponent) InProgress(user slack.User) {
	repository := c.SelectedRepository()
	branch := c.SelectedBranch()
	target := c.SelectedTarget()

	c.attachments = []slack.Attachment {
		{
			Color: "warning",
			Title: fmt.Sprintf(":hammer_and_wrench: The Make command is running ..."),
			Text: fmt.Sprintf("REPOSITORY: %s\nBRANCH: %s\nTARGET: %s", repository, branch, target),
			Footer: fmt.Sprintf("executed by %s", user.Name),
			FooterIcon: user.Profile.Image32,
		},
	}
}

func (c *MakeComponent) Done(result *models.ExecMakeResult) *MakeOutputComponent {
	if result.Success {
		c.attachments[0].Title = ":tada: make command SUCCESS!!"
		c.attachments[0].Color = "good"
	} else {
		c.attachments[0].Title = ":rain_cloud: make command FAILED..."
		c.attachments[0].Color = "error"
	}

	return &MakeOutputComponent{
		Channel:  c.Channel,
		ThreadTs: c.Ts,
		Result:   result,
	}
}


func (c *MakeComponent) Cancel(user slack.User) {
	c.attachments = []slack.Attachment {
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

func (c *MakeComponent) changeSetting(selectedAction *slack.AttachmentAction) {

	selectedOptionValue := selectedAction.SelectedOptions[0].Value

	for i := range c.attachments {
		attachment := &c.attachments[i]

		if attachment.Actions[0].Name == selectedAction.Name {
			for _, option := range attachment.Actions[0].Options {
				if option.Value == selectedOptionValue {
					// selected
					attachment.Color = "good"
					attachment.Actions = attachment.Actions[0:1] // remove cancel button
					attachment.Actions[0].SelectedOptions = []slack.AttachmentActionOption{option}
					c.attachments = c.attachments[0:i+1]
					return
				}
			}
		}
	}
}

func (c *MakeComponent) newCancelAction() slack.AttachmentAction {
	return slack.AttachmentAction{
		Name:  ActionCancel,
		Text:  "Cancel",
		Type:  "button",
		Style: "danger",
	}
}

