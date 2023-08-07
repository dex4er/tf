package run

import (
	"os"
	"os/exec"

	"github.com/dex4er/tf/util"
)

func Init(args []string) error {
	patternIgnoreLine := `Finding .* versions matching` +
		`|Finding latest version of` +
		`|from the shared cache directory` +
		`|https://www\.terraform\.io/docs/cli/plugins/signing\.html` +
		`|If you'd like to know more about provider signing, you can read about it here:` +
		`|in \.terraform/modules/` +
		`|in the \.terraform\.lock\.hcl file. Review those changes and commit them to your` +
		`|Initializing (modules` +
		`|Initializing Terraform` +
		`|Partner and community providers are signed by their developers\.` +
		`|provider plugins)\.\.\.` +
		`|Reusing previous version of` +
		`|Terraform has made some changes to the provider dependency selections recorded` +
		`|the backend` +
		`|Upgrading modules\.\.\.` +
		`|Using previously-installed` +
		`|version control system if they represent changes you intended to make\.`

	patternIgnoreFooter := `Terraform.* has been successfully initialized!`

	newArgs := []string{}

	var codesign = false

	for _, arg := range args {
		switch arg {
		case "-codesign":
			codesign = true
		default:
			newArgs = append(newArgs, arg)
		}
	}

	if err := terraformWithFilter("init", newArgs, patternIgnoreLine, patternIgnoreFooter); err != nil {
		return err
	}

	if codesign {
		if err := runCodesign(); err != nil {
			return err
		}
	}

	return nil
}

func runCodesign() error {
	cmd := exec.Command("sh", "-c", "find .terraform/providers -type f -follow -name '*_v*.*.*' | xargs -n1 codesign --force --deep --sign -")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := util.ReplacePatternInFile(".terraform.lock.hcl", `(?s)  hashes = \[.*?  \]\r?\n`, ""); err != nil {
		return err
	}

	return nil
}
