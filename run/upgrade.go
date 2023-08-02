package run

func Upgrade(args []string) error {
	return Init(append(args, "-upgrade"))
}
