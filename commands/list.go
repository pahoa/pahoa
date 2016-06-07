package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/pahoa/pahoa/core"
)

var listCmd = &cobra.Command{
	Use:  "list",
	RunE: listCmdRun,
}

func init() {
	initClientCommand(listCmd)
}

func listCmdRun(cmd *cobra.Command, args []string) error {
	url := clientOptions.endpoint + "/cards"
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var cards []core.Card

	if err := json.NewDecoder(res.Body).Decode(&cards); err != nil {
		log.Fatal("Unable to decode list of cards")
	}

	if len(cards) == 0 {
		fmt.Println("Not found any cards")
		return nil
	}

	for _, card := range cards {
		fmt.Printf("- Card: %s - Step: %s\n", card.ExternalID, card.CurrentStep)
	}

	return nil
}
