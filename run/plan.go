package run

func Plan(args []string) error {
	return commandWithProgress("plan", args)
}
