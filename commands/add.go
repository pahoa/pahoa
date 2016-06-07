package commands

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/pahoa/pahoa/core"
)

var addCardOptions = &core.AddCardOptions{}

var addCmd = &cobra.Command{
	Use:  "add",
	RunE: addCmdRun,
}

func init() {
	initClientCommand(addCmd)

	addCmd.PersistentFlags().StringVar(
		&addCardOptions.ExternalID,
		"external-id",
		"",
		"external id")
	addCmd.PersistentFlags().StringVar(
		&addCardOptions.PreviousStep,
		"previous-step",
		"",
		"card's previous step")
	addCmd.PersistentFlags().StringVar(
		&addCardOptions.CurrentStep,
		"current-step",
		"",
		"card's current step")
	addCmd.PersistentFlags().BoolVar(
		&addCardOptions.BypassActions,
		"bypass-actions",
		false,
		"ignore actions that should be executed after transition")
}

func addCmdRun(cmd *cobra.Command, args []string) error {
	data, err := json.Marshal(addCardOptions)
	if err != nil {
		return err
	}

	url := clientOptions.endpoint + "/cards"
	res, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	if res.StatusCode != 201 {
		body, _ := ioutil.ReadAll(res.Body)
		log.Fatalf("Failed to add card [%d]: %v", res.StatusCode, string(body))
	}

	defer res.Body.Close()

	var card core.Card
	if err := json.NewDecoder(res.Body).Decode(&card); err != nil {
		return err
	}

	log.Printf("Added Card: %s", card.ExternalID)

	return nil
}
