package run

func commandWithProgress(command string, args []string) error {
	ignoreLinePattern := "Terraform used the selected providers to generate the following execution" +
		"|Preparing the remote plan\\.\\.\\." +
		"|Running plan in Terraform Cloud\\. Output will stream here\\. Pressing Ctrl-C" +
		"|will stop streaming the logs, but will not stop the plan running remotely\\." +
		"|The remote workspace is configured to work with configuration at" +
		"|relative to the target repository." +
		"|excluding files or directories as defined by a \\.terraformignore file" +
		"|at .*\\.terraformignore (if it is present)," +
		"|in order to capture the filesystem context the remote workspace expects:" +
		"|plan. Resource actions are indicated with the following symbols:" +
		"|Waiting for the plan to start\\.\\.\\." +
		"|Terraform will perform the following actions:" +
		"|Terraform will perform the actions described above\\." +
		"|Terraform has compared your real infrastructure against your configuration" +
		"|and found no differences, so no changes are needed\\." +
		"|Unless you have made equivalent changes to your configuration, or ignored the" +
		"|relevant attributes using ignore_changes, the following plan may include" +
		"|actions to undo or respond to these changes\\." +
		"|This is a refresh-only plan, so Terraform will not take any actions to undo" +
		"|these\\. If you were expecting these changes then you can apply this plan to" +
		"|Terraform has checked that the real remote objects still match the result of" +
		"|your most recent changes, and found no differences\\." +
		"|To perform exactly these actions, run the following command to apply:" +
		"|To see the full warning notes, run Terraform without -compact-warnings\\." +
		"|Acquiring state lock\\. This may take a few moments\\.\\.\\." +
		"|Releasing state lock\\. This may take a few moments\\.\\.\\." +
		"|Warnings:" +
		"|Note: You .* use the -out option to save this plan, so Terraform" +
		"|guarantee to take exactly these actions if you run \"terraform apply\" now\\." +
		"|Apply complete! Resources: 0 added, 0 changed, 0 destroyed\\." +
		"|─────────────────────────────────────────────────────────────────────────────"

	// not full
	ignoreLinePattern += "|= \\(known after apply\\)" +
		"|\\(\\d+ unchanged \\w+ hidden\\)" +
		"|\\(config refers to values not yet known\\)"

	// compact
	ignoreLinePattern += "|^\\s\\s[\\s+~-]" +
		"|\\(config refers to values not yet known\\)"

	// refreshing
	ignoreLinePattern += "|: Refreshing state\\.\\.\\." +
		"|: Refreshing\\.\\.\\." +
		"|: Drift detected"

	// reading
	ignoreLinePattern += "|: Reading\\.\\.\\." +
		"|: Read complete after"

	// creating
	ignoreLinePattern += "|: Creating\\.\\.\\." +
		"|: Creation complete after"

	// destroying
	ignoreLinePattern += "|: Destroying\\.\\.\\." +
		"|: Destruction complete after"

	// modifying
	ignoreLinePattern += "|: Modifying\\.\\.\\." +
		"|: Modifications complete after"

	// still...
	ignoreLinePattern += "|: Still .*ing\\.\\.\\."

	ignoreNextLinePattern := "record the updated values in the Terraform state without changing any remote" +
		"|record the updated values in the Terraform state without changing any remote" +
		"| Experimental feature .* is active" +
		"|Saved the plan to: terraform\\.tfplan"

	ignoreBlockStartPattern := "Warning:.*Applied changes may be incomplete" +
		"|Warning:.*Resource targeting is in effect" +
		"|This plan was saved to: terraform.tfplan"

	ignoreBlockEndPattern := "suggests to use it as part of an error message" +
		"|exceptional situations such as recovering from errors or mistakes" +
		"|terraform apply \"terraform\\.tfplan\""

	footerPattern := ""

	return commandWithFilter(command, args, ignoreLinePattern, ignoreNextLinePattern, ignoreBlockStartPattern, ignoreBlockEndPattern, footerPattern)
}
