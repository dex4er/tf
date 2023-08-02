package run

func commandWithProgress(command string, args []string) error {
	ignoreLinePattern := "\\(known after apply\\)"

	ignoreNextLinePattern := "record the updated values in the Terraform state without changing any remote"

	ignoreBlockStartPattern := "Warning:.*Applied changes may be incomplete" +
		"|Warning:.*Resource targeting is in effect" +
		"|This plan was saved to: terraform.tfplan"

	ignoreBlockEndPattern := "suggests to use it as part of an error message" +
		"|exceptional situations such as recovering from errors or mistakes" +
		"|terraform apply \"terraform\\.tfplan\""

	footerPattern := ""

	return commandWithFilter(command, args, ignoreLinePattern, ignoreNextLinePattern, ignoreBlockStartPattern, ignoreBlockEndPattern, footerPattern)
}
