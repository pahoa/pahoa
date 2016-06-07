package core

import "log"

const (
	StatusMergeRequestToDevelop = "MergeRequestToDevelop"
)

type Action func()

var actions map[string]Action = map[string]Action{
	StatusMergeRequestToDevelop: mergeRequestToDevelop,
}

func mergeRequestToDevelop() {
	log.Print("MergeRequestToDevelop")
}

func CardAction(key string) Action {
	return actions[key]
}
