package run

func Rm(args []string) error {
	return terraformWithArgsQuotes("state", append([]string{"rm"}, args...))
}
