package alertsink

import (
	"os"
	"strings"

	"github.com/mobalyticshq/alertsforge/config"
	"github.com/mobalyticshq/alertsforge/sharedtools"
	"github.com/slack-go/slack"
)

type SlackSink struct {
	runbooks    *config.RunbooksConfig
	SlackClient *slack.Client
}

func NewSlackSink(runbooks *config.RunbooksConfig) *SlackSink {
	slackSink := &SlackSink{
		runbooks:    runbooks,
		SlackClient: slack.New(os.Getenv("AF_SLACK_TOKEN")),
	}

	return slackSink
}

func (o SlackSink) SendAlerts(alerts []sharedtools.Alert) (accepted []string, resolved []string, errors []error) {
	for _, alert := range alerts {

		channel, err := sharedtools.TemplateString(o.runbooks.SlackMessage.Channel, alert)
		if err != nil {
			errors = append(errors, err)
			accepted = append(accepted, alert.Fingerprint)
			continue
		}

		text := sharedtools.MustTemplateString(o.runbooks.SlackMessage.Message, alert, "error while parsing text message")

		channels := strings.Split(channel, ",")
		deduplicatedChannels := sharedtools.RemoveDuplicateFromSlice(channels)
		for _, channel := range deduplicatedChannels {
			_, _, _, err = o.SlackClient.SendMessage(channel, slack.MsgOptionText(text, false))
			if err != nil {
				errors = append(errors, err)
			}
		}
		if err == nil {
			if alert.Status == sharedtools.Resolved {
				resolved = append(resolved, alert.Fingerprint)
			} else {
				accepted = append(accepted, alert.Fingerprint)
			}
		}
	}

	return
}
