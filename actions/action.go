package actions

type Action func()

var actions map[string]Action = make(map[string]Action)

func FromName(name string) Action {
	return actions[name]
}
