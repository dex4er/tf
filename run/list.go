package run

import "github.com/dex4er/tf/util"

func List(args []string) error {
	newArgs := []string{"list"}

	for _, arg := range args {
		if util.StartsWith(arg, '-') {
			newArgs = append(newArgs, arg)
		} else {
			newArgs = append(newArgs, util.AddQuotes(arg))
		}
	}

	patternIgnoreFooter := ""
	return terraformWithoutColors("state", newArgs, patternIgnoreFooter)
}
