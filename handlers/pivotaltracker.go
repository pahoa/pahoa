package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"

	"github.com/pahoa/pahoa/core"
)

func PivotalTrackerStartCard(config *viper.Viper, card *core.Card) error {
	token := config.GetString("pivotaltracker.token")

	story, err := PivotalTrackerStoryDetail(token, card.ExternalID)
	if err != nil {
		return err
	}

	err = PivotalTrackerStoryUpdate(&PivotalTrackerStoryUpdateOptions{
		Token:        token,
		StoryID:      story.ID,
		ProjectID:    story.ProjectID,
		CurrentState: PivotalTrackerStoryStateStarted,
	})
	if err != nil {
		return err
	}

	return nil
}

func init() {
	Register(core.ActionStartCard, PivotalTrackerStartCard)
}

type PivotalTrackerStory struct {
	ID           int    `json:"id"`
	ProjectID    int    `json:"project_id"`
	CurrentState string `json:"current_state"`
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
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, err
	}

	var data PivotalTrackerStory

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

const (
	PivotalTrackerStoryStateStarted = "started"
)

type PivotalTrackerStoryUpdateOptions struct {
	Token        string
	StoryID      int
	ProjectID    int
	CurrentState string
}

func PivotalTrackerStoryUpdate(opts *PivotalTrackerStoryUpdateOptions) error {
	data, err := json.Marshal(struct {
		CurrentState string `json:"current_state"`
	}{
		CurrentState: opts.CurrentState,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://www.pivotaltracker.com/services/v5/projects/%d/stories/%d", opts.ProjectID, opts.StoryID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("X-TrackerToken", opts.Token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)
	if res != nil {
		log.Printf("status code", res.StatusCode)
	}
	if err != nil || res.StatusCode != http.StatusOK {
		log.Printf("err %#v", err)
		return err
	}

	return nil
}
