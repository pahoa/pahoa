package core

type Transition struct {
	From    string
	To      string
	Limit   int
	Actions []Action
}

type Board struct {
	Transitions []*Transition
}

func (b *Board) AddTransition(transition *Transition) {
	b.Transitions = append(b.Transitions, transition)
}

func (b *Board) GetTransition(from, to string) *Transition {
	for _, t := range b.Transitions {
		if t.From != from || t.To != to {
			continue
		}

		return t
	}

	return nil
}

func (b *Board) IsTransitionValid(from, to string) bool {
	t := b.GetTransition(from, to)

	return t != nil
}
