package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pahoa/pahoa/core"
)

var listCmdConfig = viper.New()

var listCmd = &cobra.Command{
	Use:  "list",
	RunE: listCmdRun,
}

func init() {
	initClientConfig(listCmdConfig, listCmd)

	pf := listCmd.PersistentFlags()
	pf.String("step", "", "filter by step")
	listCmdConfig.BindPFlag("step", pf.Lookup("step"))
}

func listCmdRun(cmd *cobra.Command, args []string) error {
	qs := url.Values{}
	qs.Set("step", listCmdConfig.GetString("step"))

	u := fmt.Sprintf("%s/cards?%s", listCmdConfig.GetString("endpoint"), qs.Encode())

	res, err := http.Get(u)
	if err != nil || res.StatusCode != 200 {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var cards []core.Card

	if err := json.NewDecoder(res.Body).Decode(&cards); err != nil {
		log.Fatal("Unable to decode list of cards")
	}

	if len(cards) == 0 {
		fmt.Println("No cards found")
		return nil
	}

	for _, card := range cards {
		fmt.Printf("- Card: %s (%s) - Step: %s\n", card.ID, card.Status,
			card.CurrentStep)
	}

	return nil
}
