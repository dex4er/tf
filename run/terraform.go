package run

import (
	"os"
)

var terraformGlobalArgs = []string{}

func SetTerraformGlobalArgs(args []string) {
	terraformGlobalArgs = append([]string{}, args...)
}

func terraformCommandArgs(command string, args []string) []string {
	allArgs := make([]string, 0, len(terraformGlobalArgs)+1+len(args))
	allArgs = append(allArgs, terraformGlobalArgs...)
	allArgs = append(allArgs, command)
	allArgs = append(allArgs, args...)

	return allArgs
}

func Terraform(command string, args []string) error {
	cmd := execTerraformCommand(terraformCommandArgs(command, args)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
