package core

import (
	"errors"
	"log"
)

var (
	ErrCardStatusNotOK     = errors.New("Card status is not OK")
	ErrInvalidCardMove     = errors.New("Invalid card move")
	ErrCardNotFound        = errors.New("Card not found")
	ErrWorkInProgressLimit = errors.New("Limit of work in progress cards for the given step has reached")
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

func AddCard(b *Board, m Model, e *Executor, opts *AddCardOptions) (*Card, error) {
	transition := b.GetTransition(opts.PreviousStep, opts.CurrentStep)
	if transition == nil {
		return nil, ErrInvalidCardMove
	}

	limit := b.GetLimit(opts.CurrentStep)
	if limit > 0 {
		// the following code is not concurrent-safe
		cards, err := m.ListCards(opts.CurrentStep)
		if err != nil {
			return nil, err
		}
		if len(cards) >= limit {
			return nil, ErrWorkInProgressLimit
		}
	}

	modelOptions := &ModelAddCardOptions{
		ID:           opts.ID,
		PreviousStep: opts.PreviousStep,
		CurrentStep:  opts.CurrentStep,
		Status:       CardStatusWaiting,
	}
	if opts.BypassActions {
		modelOptions.Status = CardStatusOK
	}

	card, err := m.AddCard(modelOptions)
	if err != nil {
		return nil, err
	}

	if opts.BypassActions {
		return card, nil
	}

	e.Execute(&ExecutorTask{
		Card:    card,
		Actions: transition.Actions,
	})

	return card, nil
}

type MoveCardOptions struct {
	ID            string
	To            string
	BypassActions bool
}

func MoveCard(b *Board, m Model, e *Executor, opts *MoveCardOptions) (*Card, error) {
	log.Printf("MoveCard(_, _, _, %#v)", opts)

	card, err := m.GetCard(opts.ID)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, ErrCardNotFound
	}
	if card.Status != CardStatusOK {
		return nil, ErrCardStatusNotOK
	}

	transition := b.GetTransition(card.CurrentStep, opts.To)
	if transition == nil {
		return nil, ErrInvalidCardMove
	}

	limit := b.GetLimit(opts.To)
	if limit > 0 {
		// the following code is not concurrent-safe
		cards, err := m.ListCards(opts.To)
		if err != nil {
			return nil, err
		}
		if len(cards) >= limit {
			return nil, ErrWorkInProgressLimit
		}
	}

	modelOptions := &ModelUpdateCardOptions{
		ID:           opts.ID,
		PreviousStep: card.CurrentStep,
		CurrentStep:  opts.To,
		Status:       CardStatusWaiting,
	}
	if opts.BypassActions {
		modelOptions.Status = CardStatusOK
	}

	card, err = m.UpdateCard(modelOptions)
	if err != nil {
		return nil, err
	}

	if opts.BypassActions {
		return card, nil
	}

	e.Execute(&ExecutorTask{
		Card:    card,
		Actions: transition.Actions,
	})

	return card, nil
}

func ListCards(m Model, step string) ([]*Card, error) {
	return m.ListCards(step)
}
