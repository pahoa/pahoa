package handlers

import (
	"github.com/pahoa/pahoa/core"
	"github.com/pahoa/pahoa/handlers/pivotaltracker"
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

func init() {
	Register("pivotaltracker.StartStory", pivotaltracker.StartStory)
	Register("pivotaltracker.UnstartStory", pivotaltracker.UnstartStory)
	Register("pivotaltracker.FinishStory", pivotaltracker.FinishStory)
	Register("pivotaltracker.DeliveryStory", pivotaltracker.DeliveryStory)
	Register("pivotaltracker.AcceptStory", pivotaltracker.AcceptStory)
	Register("pivotaltracker.RejectStory", pivotaltracker.RejectStory)
}
