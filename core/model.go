package core

import (
	"errors"
	"log"
)

type Model struct {
	cards []*Card
}

func (m *Model) AddCard(opts *AddCardOptions) (*Card, error) {
	for _, card := range m.cards {
		if card.ID != opts.ID {
			continue
		}

		return nil, errors.New("Card already exists")
	}

	card := &Card{
		ID:           opts.ID,
		PreviousStep: opts.PreviousStep,
		CurrentStep:  opts.CurrentStep,
		Status:       CardStatusWaiting,
	}

	m.cards = append(m.cards, card)

	return card, nil
}

func (m *Model) ListCards() []*Card {
	return m.cards
}

func (m *Model) UpdateCardStatus(cid, status string) {
	log.Printf("UpdateCardStatus(%s, %s)", cid, status)
}

func (m *Model) ClearActionLogs(cid string) {
	log.Printf("ClearActionLog(%s)", cid)
}

func (m *Model) CreateActionLog(cid string, action Action) {
	log.Printf("CreateActionLog(%s, %v)", cid, action)
	// default status = Waiting
}

func (m *Model) UpdateActionLogStatus(cid string, action Action, status, msg string) {
	log.Printf("UpdateActionLogStatus(%s, %s, %s, %s)", cid, action, status, msg)
}
