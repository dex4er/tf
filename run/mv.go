package run

func Mv(args []string) error {
	return terraformWithArgsQuotes("state", append([]string{"mv"}, args...))
}
