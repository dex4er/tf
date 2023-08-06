package run

import "github.com/dex4er/tf/util"

func Apply(args []string) error {
	newArgs := []string{}

	for _, arg := range args {
		if util.StartsWith(arg, '-') {
			newArgs = append(newArgs, arg)
		} else {
			newArgs = append(newArgs, "-target="+util.AddQuotes(arg))
		}
	}

	return terraformWithProgress("apply", newArgs)
}
