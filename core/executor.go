package core

import (
	"fmt"

	"github.com/spf13/viper"
)

type ExecutorTask struct {
	Card    *Card
	Actions []Action
}

type Executor struct {
	model    *Model
	handlers map[Action]ActionHandler
	config   *viper.Viper
	tasks    chan *ExecutorTask
}

type NewExecutorOptions struct {
	Model    *Model
	Handlers map[Action]ActionHandler
	Config   *viper.Viper
}

func NewExecutor(opts *NewExecutorOptions) *Executor {
	return &Executor{
		model:    opts.Model,
		handlers: opts.Handlers,
		config:   opts.Config,
		tasks:    make(chan *ExecutorTask, 100),
	}
}

func (e *Executor) Start() {
	go e.Loop()
}

func (e *Executor) Execute(task *ExecutorTask) {
	cid := task.Card.ID

	e.model.ClearActionLogs(cid)

	for _, action := range task.Actions {
		e.model.CreateActionLog(cid, action)
	}

	e.model.UpdateCardStatus(cid, CardStatusWaiting)

	e.tasks <- task
}

func (e *Executor) Loop() {
	for task := range e.tasks {
		cid := task.Card.ID

		e.model.UpdateCardStatus(cid, CardStatusProcessing)

		failed := false

		for _, action := range task.Actions {
			e.model.UpdateActionLogStatus(cid, action, ActionStatusProcessing, "")

			handler := e.handlers[action]
			if handler == nil {
				msg := fmt.Sprintf("Action [%s] has no handler", action)
				e.model.UpdateActionLogStatus(cid, action, ActionStatusFailed, msg)
				failed = true
				break
			}

			if err := handler(e.config, task.Card); err != nil {
				msg := fmt.Sprintf("Failed to execute action [%s]: %#v", action, err.Error())
				e.model.UpdateActionLogStatus(cid, action, ActionStatusFailed, msg)
				failed = true
				break
			}

			e.model.UpdateActionLogStatus(cid, action, ActionStatusOK, "")
		}

		switch {
		case failed:
			e.model.UpdateCardStatus(cid, CardStatusFailed)
		case !failed:
			e.model.UpdateCardStatus(cid, CardStatusOK)
		}
	}
}
