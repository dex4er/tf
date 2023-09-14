package run

import (
	"os"
	"os/exec"
)

var TERRAFORM_PATH = os.Getenv("TERRAFORM_PATH")

func execTerraformCommand(arg ...string) *exec.Cmd {
	path := TERRAFORM_PATH
	if TERRAFORM_PATH == "" {
		path = "terraform"
	}
	return exec.Command(path, arg...)
}
