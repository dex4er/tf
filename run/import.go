package run

import (
	"strings"

	"github.com/dex4er/tf/util"
)

func Import(args []string) error {
	newArgs := []string{}
	resource := ""

	noOutputs := false

	for i, arg := range args {
		if arg == "-no-show" {
			noOutputs = true
		} else if util.StartsWith(arg, '-') {
			newArgs = append(newArgs, arg)
		} else {
			resource = util.AddQuotes(arg)
			newArgs = append(newArgs, resource)
			if len(args) > i+1 {
				id := strings.Join(args[i+1:], " ")
				newArgs = append(newArgs, id)
			}
			break
		}
	}

	if err := terraformWithProgress("import", newArgs); err != nil {
		return err
	}

	if !noOutputs {
		return Show([]string{resource})
	}

	return nil
}
