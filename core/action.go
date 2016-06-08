package core

import (
	"github.com/spf13/viper"
)

const (
	ActionStartCard                   = "StartCard"
	ActionCreateMergeRequestToDevelop = "CreateMergeRequestToDevelop"
)

type Action string

type ActionHandler func(config *viper.Viper, card *Card) error
