package actions

import (
	"log"
)

func MergeRequestToDevelop() {
	log.Print("Merge request to develop")
}

func init() {
	actions["MergeRequestToDevelop"] = MergeRequestToDevelop
}
