package run

func commandWithProgress(command string, args []string) error {
	ignorePattern := "\\(known after apply\\)"
	return commandWithFilter(command, args, ignorePattern, "")
}
