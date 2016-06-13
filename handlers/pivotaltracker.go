package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/viper"

	"github.com/pahoa/pahoa/core"
)

const (
	PivotalTrackerStoryStateUnstarted = "unstarted"
	PivotalTrackerStoryStateStarted   = "started"
)

type PivotalTrackerStoryState string

type PivotalTrackerStory struct {
	ID           int                      `json:"id"`
	ProjectID    int                      `json:"project_id"`
	CurrentState PivotalTrackerStoryState `json:"current_state"`
}

func PivotalTrackerUnstartCard(config *viper.Viper, card *core.Card) error {
	storyID, err := strconv.Atoi(card.ID)
	if err != nil {
		return err
	}

	return PivotalTrackerStoryUpdate(&PivotalTrackerStoryUpdateOptions{
		Token:        config.GetString("pivotaltracker.token"),
		StoryID:      storyID,
		CurrentState: PivotalTrackerStoryStateUnstarted,
	})
}

func PivotalTrackerStartCard(config *viper.Viper, card *core.Card) error {
	storyID, err := strconv.Atoi(card.ID)
	if err != nil {
		return err
	}

	return PivotalTrackerStoryUpdate(&PivotalTrackerStoryUpdateOptions{
		Token:        config.GetString("pivotaltracker.token"),
		StoryID:      storyID,
		CurrentState: PivotalTrackerStoryStateStarted,
	})
}

func init() {
	Register(core.ActionStartCard, PivotalTrackerStartCard)
	Register(core.ActionUnstartCard, PivotalTrackerUnstartCard)
}

func PivotalTrackerStoryDetail(token string, id string) (*PivotalTrackerStory, error) {
	url := "https://www.pivotaltracker.com/services/v5/stories/" + id
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-TrackerToken", token)

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Story id invalid or not found: %s", id)
	}

	var data PivotalTrackerStory

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

type PivotalTrackerStoryUpdateOptions struct {
	Token        string
	StoryID      int
	CurrentState PivotalTrackerStoryState
}

func PivotalTrackerStoryUpdate(opts *PivotalTrackerStoryUpdateOptions) error {
	data, err := json.Marshal(struct {
		CurrentState PivotalTrackerStoryState `json:"current_state"`
	}{
		CurrentState: opts.CurrentState,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://www.pivotaltracker.com/services/v5/stories/%d", opts.StoryID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("X-TrackerToken", opts.Token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		log.Printf("PivotalTracker error %#v ", err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		log.Printf("PivotalTracker error %#v ", string(body))
		return errors.New(string(body))
	}

	return nil
}
