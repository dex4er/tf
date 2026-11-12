package run

func Plan(args []string) error {
	return terraformWithProgress("plan", targetArgs(args, false))
}
