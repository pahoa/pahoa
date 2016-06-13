package core

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type ExecutorTask struct {
	Card    *Card
	Actions []Action
}

type Executor struct {
	model    Model
	handlers map[Action]ActionHandler
	config   *viper.Viper
	tasks    chan *ExecutorTask
}

type NewExecutorOptions struct {
	Model    Model
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

	if err := e.model.ClearActionLogs(cid); err != nil {
		log.Fatalf("ClearActionLogs(%s) error: %#v", cid, err)
	}

	for _, action := range task.Actions {
		if err := e.model.CreateActionLog(cid, action); err != nil {
			log.Fatalf("CreateActionLog(%s, %s) error: %#v", cid, action, err)
		}
	}

	if err := e.model.UpdateCardStatus(cid, CardStatusWaiting); err != nil {
		log.Fatalf("UpdateCardStatus(%s, %s) error: %#v", cid, CardStatusWaiting,
			err)
	}

	e.tasks <- task
}

func (e *Executor) Loop() {
	for task := range e.tasks {
		cid := task.Card.ID

		updateCardStatus(e.model, cid, CardStatusProcessing)

		failed := false

		for _, action := range task.Actions {
			updateActionLogStatus(e.model, cid, action, ActionStatusProcessing, "")

			handler := e.handlers[action]
			if handler == nil {
				msg := fmt.Sprintf("Action [%s] has no handler", action)
				updateActionLogStatus(e.model, cid, action, ActionStatusFailed,
					msg)

				failed = true
				break
			}

			if err := handler(e.config, task.Card); err != nil {
				msg := fmt.Sprintf("Failed to execute action [%s]: %#v", action, err.Error())
				updateActionLogStatus(e.model, cid, action, ActionStatusFailed,
					msg)

				failed = true
				break
			}

			updateActionLogStatus(e.model, cid, action, ActionStatusOK, "")
		}

		switch {
		case failed:
			updateCardStatus(e.model, cid, CardStatusFailed)
		case !failed:
			updateCardStatus(e.model, cid, CardStatusOK)
		}
	}
}

func updateCardStatus(model Model, cid, status string) {
	if err := model.UpdateCardStatus(cid, status); err != nil {
		log.Fatalf("UpdateCardStatus(%s, %s) error: %#v", cid,
			status, err)
	}
}

func updateActionLogStatus(model Model, cid string, action Action, status, msg string) {
	if err := model.UpdateActionLogStatus(cid, action, status, msg); err != nil {
		log.Fatalf("UpdateActionLogStatus(%s, %s, %s, %s) error: %#v", cid, action,
			status, msg, err)
	}
}
