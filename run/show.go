package run

import (
	"fmt"

	"github.com/dex4er/tf/util"
)

func Show(args []string) error {
	resources := []string{}
	newArgs := []string{}

	for _, arg := range args {
		if util.StartsWith(arg, '-') {
			newArgs = append(newArgs, arg)
		} else {
			resources = append(resources, util.AddQuotes(arg))
		}
	}

	if len(resources) > 0 {
		newArgs = append([]string{"show"}, newArgs...)
		for i, r := range resources {
			if i > 0 {
				fmt.Println()
			}
			if err := terraformWithoutColors("state", append(newArgs, r)); err != nil {
				return err
			}
		}
		return nil
	}

	return terraformWithoutColors("show", newArgs)
}
