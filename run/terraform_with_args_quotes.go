package run

import (
	"strings"

	"github.com/dex4er/tf/util"
)

func terraformWithArgsQuotes(command string, args []string) error {
	newArgs := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			newArgs = append(newArgs, arg)
		} else {
			newArgs = append(newArgs, util.AddQuotes(arg))
		}
	}

	// it should print a message that command expects exactly one argument
	return Terraform(command, newArgs)
}
