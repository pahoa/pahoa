package handlers

import (
	"github.com/pahoa/pahoa/core"
)

var handlers map[core.Action]core.ActionHandler = make(map[core.Action]core.ActionHandler)

func Handler(action core.Action) core.ActionHandler {
	return handlers[action]
}

func Register(action core.Action, handler core.ActionHandler) {
	handlers[action] = handler
}

func GetHandlers() map[core.Action]core.ActionHandler {
	return handlers
}
