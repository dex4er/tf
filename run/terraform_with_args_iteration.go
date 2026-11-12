package run

import (
	"fmt"
)

func terraformWithArgsIteration(command string, args []string) error {
	options, resources := splitAddresses(args)

	if len(resources) > 0 {
		for i, r := range resources {
			if i > 0 {
				fmt.Println()
			}
			if err := Terraform(command, append(options, r)); err != nil {
				return err
			}
		}
		return nil
	}

	// it should print a message that command expects exactly one argument
	return Terraform(command, options)
}
