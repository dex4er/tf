package run

func Apply(args []string) error {
	return commandWithProgress("apply", args)
}
