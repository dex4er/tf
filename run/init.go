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
		`|If you'd like to know more about provider signing, you can read about it here:` +
		`|in \.terraform/modules/` +
		`|in the \.terraform\.lock\.hcl file. Review those changes and commit them to your` +
		`|Initializing (modules` +
		`|Initializing (Terraform|OpenTofu)` +
		`|Partner and community providers are signed by their developers\.` +
		`|provider plugins)\.\.\.` +
		`|Providers are signed by their developers\.` +
		`|Reusing previous version of` +
		`|selections it made above\. Include this file in your version control repository` +
		`|signing\.html` +
		`|so that (Terraform|OpenTofu) can guarantee to make the same selections by default when` +
		`|(Terraform|OpenTofu) has created a lock file .* to record the provider` +
		`|(Terraform|OpenTofu) has made some changes to the provider dependency selections recorded` +
		`|the backend` +
		`|Upgrading modules\.\.\.` +
		`|Using previously-installed` +
		`|version control system if they represent changes you intended to make\.` +
		`|you run "(terraform|tofu) init" in the future.`

	patternIgnoreFooter := `(Terraform|OpenTofu).* has been successfully initialized!`

	newArgs := []string{}

	var codesign = false

	for _, arg := range args {
		switch util.ReplaceFirstTwoDashes(arg) {
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
