package core

import (
	"errors"
	"log"

	"github.com/pahoa/pahoa/actions"
)

var (
	InvalidCardMove = errors.New("Invalid card move")
)

type Transition struct {
	From string
	To   string
}

type Board struct {
	transitions map[Transition][]actions.Action
}

func NewBoard() *Board {
	return &Board{
		transitions: make(map[Transition][]actions.Action),
	}
}

func (b *Board) AddTransition(from, to string, transitionActions ...actions.Action) {
	b.transitions[Transition{from, to}] = transitionActions
}

func (b *Board) IsTransitionValid(from, to string) bool {
	_, valid := b.transitions[Transition{from, to}]

	return valid
}

type Card struct {
	ExternalID   string `json:"external_id"`
	PreviousStep string `json:"previous_step"`
	CurrentStep  string `json:"current_step"`
	Status       string `json:"status"`
}

type AddCardOptions struct {
	ExternalID    string `json:"external_id"`
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
	board *Board
	model *Model
	cards chan *Card
}

func NewCardActionsRunner(board *Board, model *Model) *CardActionsRunner {
	return &CardActionsRunner{
		board: board,
		model: model,
		cards: make(chan *Card, 100),
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

		card.Status = "waiting"

		from := card.PreviousStep
		to := card.CurrentStep

		log.Printf("Moving from %s to %s", from, to)
		transitionActions := c.board.transitions[Transition{from, to}]
		if transitionActions == nil {
			continue
		}

		for _, action := range transitionActions {
			action()
		}
	}
}
