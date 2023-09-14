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

	noOutputs := true

	for _, arg := range args {
		if arg == "-no-output" || arg == "-no-outputs" {
			noOutputs = true
		} else if arg == "-no-output=false" || arg == "-no-outputs=false" {
			noOutputs = false
		} else if strings.HasPrefix(arg, "-") {
			newArgs = append(newArgs, arg)
		} else if _, err := os.Stat(arg); err == nil {
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
