package run

import (
	"fmt"
	"strings"

	"github.com/dex4er/tf/util"
)

func terraformWithArgsIteration(command string, args []string) error {
	resources := []string{}
	newArgs := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			newArgs = append(newArgs, arg)
		} else {
			resources = append(resources, util.AddQuotes(arg))
		}
	}

	if len(resources) > 0 {
		for i, r := range resources {
			if i > 0 {
				fmt.Println()
			}
			if err := Terraform(command, append(newArgs, r)); err != nil {
				return err
			}
		}
		return nil
	}

	// it should print a message that command expects exactly one argument
	return Terraform(command, newArgs)
}
