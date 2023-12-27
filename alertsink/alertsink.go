package alertsink

import (
	"github.com/mobalyticshq/alertsforge/config"
	"github.com/mobalyticshq/alertsforge/sharedtools"
)

type SinkInterface interface {
	SendAlerts(alerts []sharedtools.Alert) (accepted []string, resolved []string, errors []error)
}

const (
	Oncall = "oncall"
	Slack  = "slack"
)

func NewAlertSink(sinkName string, runbooks *config.RunbooksConfig) SinkInterface {
	switch sinkName {
	case Oncall:
		return NewOncallSink(runbooks)
	case Slack:
		return NewSlackSink(runbooks)
	}
	return NewOncallSink(runbooks)
}
