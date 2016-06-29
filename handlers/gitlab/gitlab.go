package gitlab

import (
	"log"

	"github.com/spf13/viper"

	"github.com/pahoa/pahoa/core"
)

func CreateMergeRequestToDevelop(config *viper.Viper, card *core.Card) error {
	log.Printf("gitlab.CreateMergeRequestToDevelop(_, %v)", card)
	/*
		create a merge request from branch to develop for every branch related
		to the card ID. ignore in case a merge request already exists.
	*/
	return nil
}

func AcceptMergeRequestToDevelop(config *viper.Viper, card *core.Card) error {
	log.Printf("gitlab.AcceptMergeRequestToDevelop(_, %v)", card)
	/*
		accept and merge all open merge requests to develop for every branch
		related to the card ID.
	*/
	return nil
}

func CloseMergeRequestToDevelop(config *viper.Viper, card *core.Card) error {
	log.Printf("gitlab.CloseMergeRequestToDevelop(_, %v)", card)
	/*
		close any open merge request to develop for every branch related to
		the card ID.
	*/
	return nil
}

func CreateMergeRequestToQA(config *viper.Viper, card *core.Card) error {
	log.Printf("gitlab.CreateMergeRequestToQA(_, %v)", card)
	/*
		create a merge request from branch to qa for every branch related
		to the card ID. ignore in case a merge request already exists.
	*/
	return nil
}

func AcceptMergeRequestToQA(config *viper.Viper, card *core.Card) error {
	log.Printf("gitlab.AcceptMergeRequestToQA(_, %v)", card)
	/*
		accept and merge all open merge requests to qa for every branch
		related to the card ID.
	*/
	return nil
}

func CloseMergeRequestToQA(config *viper.Viper, card *core.Card) error {
	log.Printf("gitlab.CloseMergeRequestToQA(_, %v)", card)
	/*
		close any open merge request to qa for every branch related to
		the card ID.
	*/
	return nil
}

func CreateAndAcceptMergeRequestToMaster(config *viper.Viper, card *core.Card) error {
	log.Printf("gitlab.CreateAndAcceptMergeRequestToMaster(_, %v)", card)
	return nil
}

func RemoveBranches(config *viper.Viper, card *core.Card) error {
	log.Printf("gitlab.RemoveBranches(_, %v)", card)
	return nil
}

func RevertMergeRequestToQA(config *viper.Viper, card *core.Card) error {
	log.Printf("gitlab.RevertMergeRequestToQA(_, %v)", card)
	return nil
}
