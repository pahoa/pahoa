package core

import (
	"errors"
)

var (
	ErrInvalidCardMove = errors.New("Invalid card move")
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

func ListCards(m Model) ([]*Card, error) {
	return m.ListCards()
}
