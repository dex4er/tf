package run

import (
	"os"
)

func Terraform(command string, args []string) error {
	cmd := execTerraformCommand(append([]string{command}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
