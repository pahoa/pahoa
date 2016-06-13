package core

import (
	"errors"
)

var (
	ErrCardAlreadyExists = errors.New("Card already exists")
)

type ModelAddCardOptions struct {
	ID           string
	PreviousStep string
	CurrentStep  string
	Status       string
}

type Model interface {
	GetCard(id string) (*Card, error)
	AddCard(*ModelAddCardOptions) (*Card, error)
	ListCards() ([]*Card, error)
	UpdateCardStatus(id, status string) error
	ClearActionLogs(id string) error
	CreateActionLog(id string, action Action) error
	UpdateActionLogStatus(id string, action Action, status, msg string) error
}
