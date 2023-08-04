package run

func Destroy(args []string) error {
	return terraformWithProgress("destroy", args)
}
