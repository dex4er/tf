package run

import "strings"

func Refresh(args []string) error {
	newArgs := []string{}
	resources := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			newArgs = append(newArgs, arg)
		} else {
			resources = append(resources, arg)
		}
	}

	if len(resources) > 0 {
		return Apply(append([]string{"-refresh-only", "-auto-approve"}, args...))
	}

	return terraformWithProgress("refresh", newArgs)
}
