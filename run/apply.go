package run

func Apply(args []string) error {
	return terraformWithProgress("apply", args)
}
