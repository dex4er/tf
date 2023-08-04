package run

func List(args []string) error {
	return terraformStateWithoutColors("list", args)
}
