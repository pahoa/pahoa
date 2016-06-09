package core

import "errors"

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
