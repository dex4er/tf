package run

func Untaint(args []string) error {
	return terraformWithArgsIteration("untaint", args)
}
