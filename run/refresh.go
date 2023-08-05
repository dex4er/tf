package run

import "github.com/dex4er/tf/util"

func Refresh(args []string) error {
	newArgs := []string{}

	for _, arg := range args {
		if util.StartsWith(arg, '-') {
			newArgs = append(newArgs, arg)
		} else {
			newArgs = append(newArgs, util.AddQuotes(arg))
		}
	}

	return terraformWithProgress("refresh", newArgs)
}
