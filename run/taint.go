package run

func Taint(args []string) error {
	return terraformWithArgsIteration("taint", args)
}
