package core

import (
	"errors"
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

type Model struct {
	cards []*Card
}

type AddCardOptions struct {
	ExternalID    string `json:"external_id"`
	PreviousStep  string `json:"previous_step"`
	CurrentStep   string `json:"current_step"`
	BypassActions bool   `json:"bypass_actions"`
}

func (m *Model) AddCard(opts *AddCardOptions) (*Card, error) {
	for _, card := range m.cards {
		if card.ExternalID != opts.ExternalID {
			continue
		}

		return nil, errors.New("Card already exists")
	}

	card := &Card{
		ExternalID:   opts.ExternalID,
		PreviousStep: opts.PreviousStep,
		CurrentStep:  opts.CurrentStep,
	}

	m.cards = append(m.cards, card)

	return card, nil
}

func (m *Model) ListCards() []*Card {
	return m.cards
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
