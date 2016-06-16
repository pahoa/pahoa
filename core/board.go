package core

type Transition struct {
	From    string
	To      string
	Actions []Action
}

type Board struct {
	Limits      map[string]int
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

func (b *Board) GetLimit(step string) int {
	limit, ok := b.Limits[step]
	if !ok {
		return 0
	}
	return limit
}

func (b *Board) IsTransitionValid(from, to string) bool {
	t := b.GetTransition(from, to)

	return t != nil
}
