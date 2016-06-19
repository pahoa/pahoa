package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/pahoa/pahoa/core"
)

var (
	moveCardID            string
	moveCardTo            string
	moveCardBypassActions bool
)

var moveCmd = &cobra.Command{
	Use:  "move",
	RunE: moveCmdRun,
}

func init() {
	initClientCommand(moveCmd)

	pf := moveCmd.PersistentFlags()

	pf.StringVar(&moveCardID, "id", "", "id")
	pf.StringVar(&moveCardTo, "to", "", "to step")
	pf.BoolVar(&moveCardBypassActions, "bypass-actions", false,
		"ignore actions that should be executed after transition")
}

func moveCmdRun(cmd *cobra.Command, args []string) error {
	qs := url.Values{}
	qs.Set("bypass-actions", strconv.FormatBool(moveCardBypassActions))

	u := fmt.Sprintf("%s/cards/%s/step/%s?%s", clientOptions.endpoint, moveCardID,
		moveCardTo, qs.Encode())

	fmt.Printf("url", u)
	res, err := http.Post(u, "application/json", nil)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		log.Fatalf("Failed to move card [%d]: %v", res.StatusCode, string(body))
	}

	defer res.Body.Close()

	var card core.Card
	if err := json.NewDecoder(res.Body).Decode(&card); err != nil {
		return err
	}

	log.Printf("Moved Card: %s", card.ID)

	return nil
}
