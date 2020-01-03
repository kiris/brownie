package brownie

import (
	"fmt"
	"github.com/nlopes/slack"
	"strings"
)

type ExecMakeResponse struct {
	event *slack.MessageEvent
	result *MakeResult
}

func (m *ExecMakeResponse) Send(rtm *slack.RTM) error {
	channel, ts, err := m.sendResultMessage(rtm)
	if err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	if _, _, err := m.sentResultDetailMessage(rtm, channel, ts); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}

func (m *ExecMakeResponse) sendResultMessage(rtm *slack.RTM) (string, string, error) {
	user, err := rtm.GetUserInfo(m.event.User)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user info: %s", m.event.User)
	}

	attachment := slack.MsgOptionAttachments(
		slack.Attachment{
			Title: m.title(),
			Color: m.color(),
			Fields: []slack.AttachmentField{
				{
					Title: "project",
					Value: m.result.project.name,
				},
				{
					Title: "branch",
					Value: m.result.branch,
				},
				{
					Title: "targets",
					Value: m.targets(),
				},
			},
			Footer: fmt.Sprintf("Executed by %s", user.Name),
			FooterIcon: user.Profile.Image32,
		},
	)

	return rtm.PostMessage(m.event.Channel, attachment)
}

func (m *ExecMakeResponse) title() string {
	if m.result.success {
		return ":tada: make SUCCESS!!"
	} else {
		return ":rain_cloud: make FAILED..."
	}
}

func (m *ExecMakeResponse) color() string {
	if m.result.success {
		return "good"
	} else {
		return "danger"
	}
}

func (m *ExecMakeResponse) targets() string {
	if len(m.result.targets) == 0 {
		return "(default)"
	} else {
		return strings.Join(m.result.targets, " ")
	}
}

func (m *ExecMakeResponse) sentResultDetailMessage(rtm *slack.RTM, channel string, timestamp string) (string, string, error) {
	ts := slack.MsgOptionTS(timestamp)
	attachment := slack.MsgOptionAttachments(
		slack.Attachment{
			Title: ":memo: more details",
			Color: "none",
			Fields: []slack.AttachmentField{
				{
					Title: "exec command",
					Value: m.result.exec,
				},
				{
					Title: "output",
					Value: m.result.output,
				},
			},
		},
	)

	return rtm.PostMessage(channel, ts, attachment)
}

