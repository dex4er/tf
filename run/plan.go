package run

func Plan(args []string) error {
	return terraformWithProgress("plan", args)
}
