package run

func Destroy(args []string) error {
	return commandWithProgress("destroy", args)
}
