package run

func commandWithProgress(command string, args []string) error {
	return commandWithFilter(command, args, "\\(known after apply\\)", "")
}
