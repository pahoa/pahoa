package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pahoa/pahoa/core"
)

var addCardConfig = viper.New()

var addCmd = &cobra.Command{
	Use:     "add <id>",
	PreRunE: addCmdPreRun,
	RunE:    addCmdRun,
}

func init() {
	initClientConfig(addCardConfig, addCmd)

	pf := addCmd.PersistentFlags()

	pf.String("previous-step", "", "card's previous step")
	addCardConfig.BindPFlag("previous-step", pf.Lookup("previous-step"))

	pf.String("current-step", "todo", "card's current step")
	addCardConfig.BindPFlag("current-step", pf.Lookup("current-step"))

	pf.Bool("bypass-actions", false,
		"ignore actions that should be executed after transition")
	addCardConfig.BindPFlag("bypass-actions", pf.Lookup("bypass-actions"))
}

func addCmdPreRun(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("id is required")
	}

	return nil
}

func addCmdRun(cmd *cobra.Command, args []string) error {
	opts := &core.AddCardOptions{
		ID:            args[0],
		PreviousStep:  addCardConfig.GetString("previous-step"),
		CurrentStep:   addCardConfig.GetString("current-step"),
		BypassActions: addCardConfig.GetBool("bypass-actions"),
	}
	data, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	url := addCardConfig.GetString("endpoint") + "/cards"
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

	log.Printf("Added Card: %s", card.ID)

	return nil
}
