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
	"github.com/spf13/viper"

	"github.com/pahoa/pahoa/core"
)

var moveCmdConfig = viper.New()

var moveCmd = &cobra.Command{
	Use:     "move <id> <step>",
	PreRunE: moveCmdPreRun,
	RunE:    moveCmdRun,
}

func init() {
	initClientConfig(moveCmdConfig, moveCmd)

	pf := moveCmd.PersistentFlags()

	pf.Bool("bypass-actions", false,
		"ignore actions that should be executed after transition")
	moveCmdConfig.BindPFlag("bypass-actions", pf.Lookup("bypass-actions"))
}

func moveCmdPreRun(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("<id> and <step> are required")
	}

	return nil
}

func moveCmdRun(cmd *cobra.Command, args []string) error {
	qs := url.Values{}
	qs.Set("bypass-actions", strconv.FormatBool(moveCmdConfig.GetBool("bypass-actions")))

	u := fmt.Sprintf("%s/cards/%s/step/%s?%s", moveCmdConfig.GetString("endpoint"),
		args[0], args[1], qs.Encode())

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
