package run

import (
	"strings"

	"github.com/dex4er/tf/util"
)

func List(args []string) error {
	newArgs := []string{"list"}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			newArgs = append(newArgs, arg)
		} else {
			newArgs = append(newArgs, util.AddQuotes(arg))
		}
	}

	noOutputs := false

	return terraformWithoutColors("state", noOutputs, newArgs)
}
