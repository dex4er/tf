package run

import (
	"os"
	"os/exec"
)

func Terraform(command string, args []string) error {
	cmd := exec.Command("terraform", append([]string{command}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
