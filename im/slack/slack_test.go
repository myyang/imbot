package slack

import "testing"

func TestSlack(t *testing.T) {
	slackCli := New()

	opt := map[string]interface{}{
		"channel_id": "#bot-test",
	}
	logger := slackCli.Logger(opt)
	logger.Infof("unit-test-info")
}
