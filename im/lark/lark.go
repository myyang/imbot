package lark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	lAPI "github.com/larksuite/oapi-sdk-go/api"
	lReq "github.com/larksuite/oapi-sdk-go/api/core/request"
	lCore "github.com/larksuite/oapi-sdk-go/core"
	lCfg "github.com/larksuite/oapi-sdk-go/core/config"
	lConstants "github.com/larksuite/oapi-sdk-go/core/constants"
	lEventHdr "github.com/larksuite/oapi-sdk-go/event/core/handlers"
	lEventGin "github.com/larksuite/oapi-sdk-go/event/http/gin"

	botCmd "github.com/myyang/imbot/commands"
	botLog "github.com/myyang/imbot/log"
)

var (
	once      sync.Once
	mutex     sync.Mutex
	singleton *Lark
)

type Lark struct {
	config *lCfg.Config
}

const (
	larkPath = "/lark"

	larkEventMessage = "message"
)

// RegisterGin once with pre-defined path
func RegisterGin(r *gin.Engine) {
	once.Do(func() { sdkClient(r) })
}

func sdkClient(r *gin.Engine) {
	// SDK parse 'type' in event object, not the type: 'event_callback' at top
	// level
	lEventHdr.SetTypeHandler(
		larkConfig,
		larkEventMessage,
		&eventCallbackHandler{lark: New()},
	)
	lEventGin.Register(larkPath, larkConfig, r)
}

func customClient(r *gin.Engine) {
	// maybe for debugging and fallback.
	r.POST(larkPath, New().handleHttp)
}

// New return lark bot instance.
func New() *Lark {
	mutex.Lock()
	defer mutex.Unlock()

	if singleton != nil {
		return singleton
	}

	singleton = &Lark{
		config: larkConfig,
	}

	return singleton
}

type larkEncryptedData struct {
	Encrypt string `json:"encrypt"`
}

func (l *Lark) handleHttp(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// handle message
	msg := map[string]interface{}{}
	err = json.Unmarshal(raw, &msg)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !l.validate(msg) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	msgType, ok := msg["type"]
	if ok {
		// version 1.0
		switch msgType.(string) {
		case string(lConstants.CallbackTypeChallenge):
			c.JSON(http.StatusOK, gin.H{"challenge": msg["challenge"]})
			return
		case string(lConstants.CallbackTypeEvent):
			go l.parseEvent(msg)
		}
	} else {
		// TODO: version 2.0+
	}

	c.JSON(http.StatusOK, struct{}{})
}

func (l *Lark) validate(msg map[string]interface{}) bool {
	if msg == nil {
		return false
	}

	if token, ok := msg["token"]; ok && token.(string) == os.Getenv("LARK_MSG_TOKEN") {
		return true
	}

	// message version 2.0+
	if schema, ok := msg["schema"]; ok {
		switch schema.(string) {
		case "2.0":
			if header, ok := msg["header"]; ok {
				return l.validate(header.(map[string]interface{}))
			}
		}
	}

	return false
}

/*
{"event":
{
	"app_id":"cli_a0c31007c7b85013",
	"chat_type":"group",
	"employee_id":"72dagf2e",
	"is_mention":true,"lark_version":"lark/3.42.8",
	"message_id":"",
	"msg_type":"text",
	"open_chat_id":"oc_5226dc3306574fae5a2ca1dca6b1a859",
	"open_id":"ou_2988d05af705b066d80c14293649659c",
	"open_message_id":"om_26021f1187ce32fff8a1a144ebec437b",
	"parent_id":"",
	"root_id":"",
	"tenant_key":"2e94517458cf1652",
	"text":"\u003cat open_id=\"ou_8b3ce8e1601a68ca9b9711aa75935917\"\u003e@pgbot\u003c/at\u003e show me",
	"text_without_at_bot":" show me",
	"type":"message",
	"union_id":"on_1760006daff55dd346feaa3f89facc24",
	"user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36 Lark/3.42.8 LarkLocale/en_US SDK-Version/3.42.22",
	"user_open_id":"ou_2988d05af705b066d80c14293649659c"
}
"token":"qwRsWcYkg4LAoeuJjhVPmhs6rVGNPBUh",
"ts":"1616480799.173236",
"type":"event_callback",
"uuid":"1d6a5257a693fc849b853a8bf774bd31"
}
*/
func (l *Lark) parseEvent(msg map[string]interface{}) {
	event := msg["event"].(map[string]interface{})
	if is_mention, ok := event["is_mention"]; !ok || !is_mention.(bool) {
		// skip not mentioned msg
		return
	}

	inf, ok := event["text_without_at_bot"]
	if !ok {
		fmt.Printf("no text without bot. event:\n%v\n", msg)
		return
	}

	literal, ok := inf.(string)
	if !ok {
		fmt.Printf("test is not string. event:\n%v\n", msg)
		return
	}

	chatID := event["open_chat_id"].(string)
	openID := event["open_id"].(string)
	logOpt := map[string]interface{}{
		"chat_id": chatID,
		"open_id": openID,
		"root_id": event["open_message_id"],
	}

	ctx := context.WithValue(
		context.Background(),
		botLog.CtxLogger,
		l.Logger(logOpt),
	)
	botCmd.Execute(ctx, strings.Split(strings.Trim(literal, " "), " "))
}

func (l *Lark) botLog(msg string) { l.botSay(os.Getenv("LARK_LOG_CHAT_ID"), msg) }

func (l *Lark) botSay(channelID string, msg string) {
	body := map[string]interface{}{
		"chat_id":  channelID,
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": msg,
		},
	}

	l.send(body)
}

func (l *Lark) botSayOpt(logOpt map[string]interface{}, msg string) {
	body := map[string]interface{}{
		"chat_id":  logOpt["chat_id"],
		"root_id":  logOpt["root_id"],
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": msg,
		},
	}

	l.send(body)
}

func (l *Lark) send(body map[string]interface{}) {
	ret := map[string]interface{}{}
	req := lReq.NewRequestWithNative(
		"message/v4/send",
		http.MethodPost,
		lReq.AccessTokenTypeTenant,
		body,
		&ret,
	)

	ctx := lCore.WrapContext(context.Background())
	err := lAPI.Send(ctx, l.config, req)
	if err != nil {
		return
	}
}

// Logger returns a logger that wraps lark instance and set logging target to
// lark chat.
func (l *Lark) Logger(logOpt map[string]interface{}) *LarkLogger {
	return &LarkLogger{l: l, opt: logOpt}
}

// LarkLogger wraps lark instance and set logging target to lark chat.
type LarkLogger struct {
	l   *Lark
	opt map[string]interface{}
}

// Debugf adapts fmt.Printf(string, ...interface{}) and log to chat which is
// defined by LARK_LOG_CHAT_ID env.
func (l *LarkLogger) Debugf(tmpl string, args ...interface{}) {
	l.l.botLog(
		fmt.Sprintf("request from uid: %v\n", l.opt["open_id"]) +
			fmt.Sprintf(tmpl, args...),
	)
}

// Infof adapts fmt.Printf(string, ...interface{}) and log to chat which the
// message is comming from and @ the sending user.
func (l *LarkLogger) Infof(tmpl string, args ...interface{}) {
	l.l.botSayOpt(
		l.opt,
		fmt.Sprintf("<at open_id=\"%v\">user</at>\n", l.opt["open_id"])+
			fmt.Sprintf(tmpl, args...),
	)
}

// Write log to chat which the message is comming from and @ the sending user by
// sending bytes message.
func (l *LarkLogger) Write(msg []byte) (int, error) {
	l.l.botSay(l.opt["chat_id"].(string), string(msg))
	return len(msg), nil
}
