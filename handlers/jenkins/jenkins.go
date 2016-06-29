package jenkins

import (
	"log"

	"github.com/spf13/viper"

	"github.com/pahoa/pahoa/core"
)

func Build(config *viper.Viper, card *core.Card) error {
	log.Printf("jenkins.Build(_, %v)", card)
	return nil
}
