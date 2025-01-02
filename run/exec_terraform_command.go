package run

import (
	"os"
	"os/exec"

	"github.com/dex4er/tf/util"
)

var TERRAFORM_PATH = os.Getenv("TERRAFORM_PATH")

func execTerraformCommand(arg ...string) *exec.Cmd {
	path := TERRAFORM_PATH
	if TERRAFORM_PATH == "" {
		versionFile := util.FindDotVersionFile()
		if versionFile == util.OpentofuVersionFile {
			path = "tofu"
		} else {
			path = "terraform"
		}
	}
	return exec.Command(path, arg...)
}
