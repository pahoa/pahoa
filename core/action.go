package core

import (
	"github.com/spf13/viper"
)

const (
	ActionStartCard             = "StartCard"
	ActionMergeRequestToDevelop = "MergeRequestToDevelop"
)

type Action string

type ActionHandler func(config *viper.Viper, card *Card) error
