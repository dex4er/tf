package run

import (
	"strings"

	"github.com/dex4er/tf/util"
)

func terraformWithArgsQuotes(command string, args []string) error {
	return Terraform(command, quotedArgs(args))
}

// Keeps the order, as the subcommand of `state` comes first.
func quotedArgs(args []string) []string {
	newArgs := []string{}

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			optionArgs, n := optionArgs(args, i)
			newArgs = append(newArgs, optionArgs...)
			i += n - 1
			continue
		}
		newArgs = append(newArgs, util.AddQuotes(args[i]))
	}

	return newArgs
}
