package run

func Init(args []string) error {
	patternIgnoreLine := "Finding .* versions matching" +
		"|Initializing Terraform" +
		"|Initializing (modules" +
		"|the backend" +
		"|provider plugins)\\.\\.\\." +
		"|Upgrading modules\\.\\.\\." +
		"|Using previously-installed" +
		"|Reusing previous version of" +
		"|from the shared cache directory" +
		"|in \\.terraform/modules/" +
		"|Finding latest version of" +
		"|Partner and community providers are signed by their developers\\." +
		"|If you'd like to know more about provider signing, you can read about it here:" +
		"|https://www\\.terraform\\.io/docs/cli/plugins/signing\\.html" +
		"|Terraform has made some changes to the provider dependency selections recorded" +
		"|in the \\.terraform\\.lock\\.hcl file. Review those changes and commit them to your" +
		"|version control system if they represent changes you intended to make\\."

	patternIgnoreFooter := "Terraform.* has been successfully initialized!"

	return commandWithFilter("init", args, patternIgnoreLine, patternIgnoreFooter)
}
