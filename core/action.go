package core

import (
	"github.com/spf13/viper"
)

const (
	ActionStartCard                   = "StartCard"
	ActionUnstartCard                 = "UnstartCard"
	ActionCreateMergeRequestToDevelop = "CreateMergeRequestToDevelop"
)

const (
	ActionStatusWaiting    = "waiting"
	ActionStatusProcessing = "processing"
	ActionStatusOK         = "ok"
	ActionStatusFailed     = "failed"
)

type Action string

type ActionHandler func(config *viper.Viper, card *Card) error
