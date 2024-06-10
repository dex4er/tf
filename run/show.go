package run

import (
	"fmt"
	"os"
	"strings"

	"github.com/dex4er/tf/util"
)

func Show(args []string) error {
	resources := []string{}
	newArgs := []string{}

	noOutputs := false

	for _, arg := range args {
		switch util.ReplaceFirstTwoDashes(arg) {
		case "-no-output":
			noOutputs = true
		case "-no-outputs":
			noOutputs = true
		case "-no-output=false":
			noOutputs = false
		case "-no-outputs=false":
			noOutputs = false
		default:
			if strings.HasPrefix(arg, "-") {
				newArgs = append(newArgs, arg)
			} else if _, err := os.Stat(arg); err == nil {
				newArgs = append(newArgs, arg)
			} else {
				resources = append(resources, util.AddQuotes(arg))
			}
		}
	}

	if len(resources) > 0 {
		newArgs = append([]string{"show"}, newArgs...)
		for i, r := range resources {
			if i > 0 {
				fmt.Println()
			}
			if err := terraformWithoutColors("state", noOutputs, append(newArgs, r)); err != nil {
				return err
			}
		}
		return nil
	}

	return terraformWithoutColors("show", noOutputs, newArgs)
}
