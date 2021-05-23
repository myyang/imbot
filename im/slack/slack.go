package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	botCmd "github.com/myyang/imbot/commands"
	botLog "github.com/myyang/imbot/log"
)

var (
	once      sync.Once
	mutex     sync.Mutex
	singleton *Slack
)

const (
	slackPath = "/slack"
)

type Slack struct {
	api           *slack.Client
	signingSecret string
}

// RegisterGin once with pre-defined path
func RegisterGin(r *gin.Engine) {
	r.POST(slackPath, New().handleHttp)
}

// New returns slack bot instance.
func New() *Slack {
	mutex.Lock()
	defer mutex.Unlock()

	if singleton != nil {
		return singleton
	}

	singleton = &Slack{
		api:           slack.New(os.Getenv("SLACK_TOKEN")), // OAuth Token
		signingSecret: os.Getenv("SLACK_SIGNING_SECRET"),   // App Sign Secret
	}
	return singleton
}

func (s *Slack) handleHttp(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	secretVerifier, err := slack.NewSecretsVerifier(c.Request.Header, s.signingSecret)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, err = secretVerifier.Write(raw)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = secretVerifier.Ensure()
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(raw),
		slackevents.OptionNoVerifyToken(),
	)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(raw, &r)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.String(http.StatusOK, r.Challenge)
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			opt := map[string]interface{}{
				"channel_id": ev.Channel,
			}
			ctx := context.WithValue(
				context.Background(),
				botLog.CtxLogger,
				s.Logger(opt),
			)
			botCmd.Execute(ctx, strings.Split(strings.Trim(ev.Text, " "), " "))
		}

		c.Status(http.StatusOK)
	}
}

// Logger returns a logger that wraps slack instance and set logging target to
// slack channel.
func (s *Slack) Logger(opt map[string]interface{}) (l *SlackLogger) {
	l = &SlackLogger{
		s:   s,
		opt: opt,
	}
	return
}

// SlackLogger wraps slack instance and set logging target to slack channel.
type SlackLogger struct {
	s   *Slack
	opt map[string]interface{}
}

// Debugf adapts fmt.Printf(string, ...interface{}) and log to channel which is
// defined by SLACK_LOG_CHANNEL env.
func (s *SlackLogger) Debugf(tmpl string, args ...interface{}) {
	_, _, err := s.s.api.PostMessage(
		os.Getenv("SLACK_LOG_CHANNEL"),
		slack.MsgOptionText(fmt.Sprintf(tmpl, args...), false),
	)
	if err != nil {
		log.Printf("Debugf error: %v\n", err)
	}
}

// Infof adapts fmt.Printf(string, ...interface{}) and log to channel which the
// message is comming from and @ the sending user.
func (s *SlackLogger) Infof(tmpl string, args ...interface{}) {
	_, _, err := s.s.api.PostMessage(
		s.opt["channel_id"].(string),
		slack.MsgOptionText(fmt.Sprintf(tmpl, args...), false),
	)
	if err != nil {
		log.Printf("Printf error: %v\n", err)
	}
}

// Write log to channel which the message is comming from and @ the sending user
// by sending bytes message.
func (s *SlackLogger) Write(msg []byte) (int, error) {
	s.Infof(string(msg))
	return len(msg), nil
}
