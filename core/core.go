package core

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

var (
	InvalidCardMove = errors.New("Invalid card move")
)

const (
	CardStatusWaiting    = "waiting"
	CardStatusProcessing = "processing"
	CardStatusOK         = "ok"
	CardStatusFailed     = "failed"
)

type Card struct {
	ID           string `json:"id"`
	PreviousStep string `json:"previous_step"`
	CurrentStep  string `json:"current_step"`
	Status       string `json:"status"`
}

type AddCardOptions struct {
	ID            string `json:"id"`
	PreviousStep  string `json:"previous_step"`
	CurrentStep   string `json:"current_step"`
	BypassActions bool   `json:"bypass_actions"`
}

func AddCard(b *Board, m *Model, r *CardActionsRunner, opts *AddCardOptions) (*Card, error) {
	if !b.IsTransitionValid(opts.PreviousStep, opts.CurrentStep) {
		return nil, InvalidCardMove
	}

	card, err := m.AddCard(opts)
	if err != nil {
		return nil, err
	}

	if opts.BypassActions {
		return card, nil
	}

	r.Add(card)

	return card, nil
}

func ListCards(m *Model) []*Card {
	return m.ListCards()
}

type CardActionsRunner struct {
	board    *Board
	model    *Model
	handlers map[Action]ActionHandler
	config   *viper.Viper
	cards    chan *Card
}

type NewCardActionsRunnerOptions struct {
	Board    *Board
	Model    *Model
	Handlers map[Action]ActionHandler
	Config   *viper.Viper
}

func NewCardActionsRunner(opts *NewCardActionsRunnerOptions) *CardActionsRunner {
	return &CardActionsRunner{
		board:    opts.Board,
		model:    opts.Model,
		handlers: opts.Handlers,
		config:   opts.Config,
		cards:    make(chan *Card, 100),
	}
}

func (c *CardActionsRunner) Start() {
	go c.Loop()
}

func (c *CardActionsRunner) Add(card *Card) {
	c.cards <- card
}

func (c *CardActionsRunner) Loop() {
	for card := range c.cards {
		log.Printf("Processing actions of card: %s", card.ExternalID)

		card.Status = CardStatusProcessing

		from := card.PreviousStep
		to := card.CurrentStep

		log.Printf("Moving from %s to %s", from, to)
		transition := c.board.GetTransition(from, to)

		for _, action := range transition.Actions {
			handler := c.handlers[action]
			if handler == nil {
				log.Printf("Action [%s] has no handler", action)
				card.Status = CardStatusFailed
				break
			}

			if err := handler(c.config, card); err != nil {
				log.Printf("Failed to execute action [%s]", action)
				card.Status = CardStatusFailed
				break
			}
		}

		card.Status = CardStatusOK
	}
}
