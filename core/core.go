package core

import (
	"errors"
	"log"
)

var (
	InvalidCardMove = errors.New("Invalid card move")
)

type Transition struct {
	From string
	To   string
}

type Board struct {
	transitions map[Transition][]Action
}

func NewBoard() *Board {
	return &Board{
		transitions: make(map[Transition][]Action),
	}
}

func (b *Board) AddTransition(from, to string, actions ...Action) {
	b.transitions[Transition{from, to}] = actions
}

func (b *Board) IsTransitionValid(from, to string) bool {
	_, valid := b.transitions[Transition{from, to}]

	return valid
}

type Action func()

type Card struct {
	ExternalID   string `json:"external_id"`
	PreviousStep string `json:"previous_step"`
	CurrentStep  string `json:"current_step"`
}

type AddCardOptions struct {
	ExternalID    string `json:"external_id"`
	PreviousStep  string `json:"previous_step"`
	CurrentStep   string `json:"current_step"`
	BypassActions bool   `json:"bypass_actions"`
}

func AddCard(b *Board, m *Model, opts *AddCardOptions) (*Card, error) {
	if !b.IsTransitionValid(opts.PreviousStep, opts.CurrentStep) {
		return nil, InvalidCardMove
	}

	return m.AddCard(opts)
}

func ListCards(m *Model) []*Card {
	return m.ListCards()
}

type CardActionsRunner struct {
	model *Model
	cards chan *Card
}

func NewCardActionsRunner(model *Model) *CardActionsRunner {
	return &CardActionsRunner{
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
	}
}
