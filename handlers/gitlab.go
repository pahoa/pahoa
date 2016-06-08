package handlers

import (
	"log"

	"github.com/spf13/viper"

	"github.com/pahoa/pahoa/core"
)

func GitlabCreateMergeRequestToDevelop(config *viper.Viper, card *core.Card) error {
	log.Print("Gitlab - create merge request to develop - card: %s", card.ExternalID)
	return nil
}

func init() {
	Register(core.ActionCreateMergeRequestToDevelop, GitlabCreateMergeRequestToDevelop)
}
