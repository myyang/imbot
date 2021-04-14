package lark

import (
	lCore "github.com/larksuite/oapi-sdk-go/core"
	lHandlers "github.com/larksuite/oapi-sdk-go/event/core/handlers"
)

type eventCallbackHandler struct {
	lark *Lark
}

// GetEvent implements github.com/lurksuite/event/core/handlers.Handler
// returns the response decoding target and pass to Handler.Handle as param e
func (h *eventCallbackHandler) GetEvent() interface{} { return &map[string]interface{}{} }

// Handle implements github.com/lurksuite/event/core/handlers.Handler
// handle function is the main processor
func (h *eventCallbackHandler) Handle(ctx *lCore.Context, e interface{}) (err error) {
	go h.lark.parseEvent(*(e.(*map[string]interface{})))
	return
}

var _ lHandlers.Handler = (*eventCallbackHandler)(nil)
