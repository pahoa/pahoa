package handlers

import (
	"log"

	"github.com/spf13/viper"

	"github.com/pahoa/pahoa/core"
)

func PivotalTrackerStartCard(config *viper.Viper, card *core.Card) error {
	log.Print("PivotalTracker start card: %#v", card)
	return nil
}

func init() {
	Register(core.ActionStartCard, PivotalTrackerStartCard)
}
