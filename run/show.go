package run

import (
	"fmt"

	"github.com/dex4er/tf/util"
)

func Show(args []string) error {
	resources := []string{}
	newArgs := []string{}

	noOutputs := false

	for _, arg := range args {
		if arg == "-no-output" || arg == "-no-outputs" {
			noOutputs = true
		} else if util.StartsWith(arg, '-') {
			newArgs = append(newArgs, arg)
		} else {
			resources = append(resources, util.AddQuotes(arg))
		}

	}

	patternIgnoreFooter := ""

	if len(resources) > 0 {
		newArgs = append([]string{"show"}, newArgs...)
		for i, r := range resources {
			if i > 0 {
				fmt.Println()
			}
			if err := terraformWithoutColors("state", append(newArgs, r), patternIgnoreFooter); err != nil {
				return err
			}
		}
		return nil
	}

	if noOutputs {
		patternIgnoreFooter = `^Outputs:$`
	}
	return terraformWithoutColors("show", newArgs, patternIgnoreFooter)
}
