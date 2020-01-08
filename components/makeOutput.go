package components

import (
	"github.com/kiris/brownie/model"
	"github.com/nlopes/slack"
)

type MakeOutputComponent struct {
	Channel string
	ThreadTs string
	Result *model.ExecMakeResult
}


func (c *MakeOutputComponent) Render() (string, string, string, []slack.Attachment) {
	attachments := []slack.Attachment{
		{
			Title: ":memo: details",
			Color: "none",
			Fields: []slack.AttachmentField{
				{
					Title: "exec command",
					Value: c.Result.Exec,
				},
				{
					Title: "output",
					Value: c.Result.Output,
				},
			},
		},
	}

	return c.Channel, "", c.ThreadTs, attachments
}