package pivotaltracker

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
	StoryStateUnstarted = "unstarted"
	StoryStateStarted   = "started"
	StoryStateFinished  = "finished"
	StoryStateDelivered = "delivered"
	StoryStateAccepted  = "accepted"
	StoryStateRejected  = "rejected"
)

type StoryState string

type Story struct {
	ID           int        `json:"id"`
	ProjectID    int        `json:"project_id"`
	CurrentState StoryState `json:"current_state"`
}

func UnstartStory(config *viper.Viper, card *core.Card) error {
	return ChangeStoryState(config, card, StoryStateUnstarted)
}

func StartStory(config *viper.Viper, card *core.Card) error {
	return ChangeStoryState(config, card, StoryStateStarted)
}

func FinishStory(config *viper.Viper, card *core.Card) error {
	return ChangeStoryState(config, card, StoryStateFinished)
}

func DeliveryStory(config *viper.Viper, card *core.Card) error {
	return ChangeStoryState(config, card, StoryStateDelivered)
}

func AcceptStory(config *viper.Viper, card *core.Card) error {
	return ChangeStoryState(config, card, StoryStateAccepted)
}

func RejectStory(config *viper.Viper, card *core.Card) error {
	return ChangeStoryState(config, card, StoryStateRejected)
}

func ChangeStoryState(config *viper.Viper, card *core.Card, state StoryState) error {
	storyID, err := strconv.Atoi(card.ID)
	if err != nil {
		return err
	}

	return StoryUpdate(&StoryUpdateOptions{
		Token:        config.GetString("pivotaltracker.token"),
		StoryID:      storyID,
		CurrentState: state,
	})
}

func StoryDetail(token string, id string) (*Story, error) {
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

	var data Story

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

type StoryUpdateOptions struct {
	Token        string
	StoryID      int
	CurrentState StoryState
}

func StoryUpdate(opts *StoryUpdateOptions) error {
	data, err := json.Marshal(struct {
		CurrentState StoryState `json:"current_state"`
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
