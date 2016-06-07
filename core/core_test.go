package core

import (
	"reflect"
	"testing"
)

func TestAddCard(t *testing.T) {
	model := &Model{}

	board := NewBoard()

	board.AddTransition("", "todo")
	board.AddTransition("todo", "in-development")

	expected := &Card{
		ExternalID:   "123456",
		PreviousStep: "",
		CurrentStep:  "todo",
	}

	actual, err := AddCard(board, model, &AddCardOptions{
		ExternalID:    expected.ExternalID,
		PreviousStep:  expected.PreviousStep,
		CurrentStep:   expected.CurrentStep,
		BypassActions: true,
	})

	if err != nil {
		t.Errorf("Unexpected error [%#v]", err)
	}

	if *actual != *expected {
		t.Errorf("Expected [%#v], got [%#v]", expected, actual)
	}
}

func TestListCard(t *testing.T) {
	m := &Model{}

	b := NewBoard()

	b.AddTransition("", "todo")

	expectedCards := []*Card{
		{
			ExternalID:   "1",
			PreviousStep: "",
			CurrentStep:  "todo",
		},
		{
			ExternalID:   "2",
			PreviousStep: "",
			CurrentStep:  "todo",
		},
	}

	for _, card := range expectedCards {
		_, err := AddCard(b, m, &AddCardOptions{
			ExternalID:    card.ExternalID,
			PreviousStep:  card.PreviousStep,
			CurrentStep:   card.CurrentStep,
			BypassActions: true,
		})
		if err != nil {
			t.Fatalf("Unexpected error [%#v]", err)
		}
	}

	actualCards := ListCards(m)

	if !reflect.DeepEqual(actualCards, expectedCards) {
		t.Errorf("Expected [%#v], got [%#v]", expectedCards, actualCards)
	}
}
