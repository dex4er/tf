package run

import (
	"os"
	"os/exec"
)

var TERRAFORM = os.Getenv("TERRAFORM")

func execTerraformCommand(arg ...string) *exec.Cmd {
	name := TERRAFORM
	if TERRAFORM == "" {
		name = "terraform"
	}
	return exec.Command(name, arg...)
}
