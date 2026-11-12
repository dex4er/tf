package run

import (
	"slices"
	"strings"
)

func Refresh(args []string) error {
	newArgs := targetArgs(args, false)

	if slices.ContainsFunc(newArgs, func(arg string) bool { return strings.HasPrefix(arg, "-target=") }) {
		return Apply(append([]string{"-refresh-only", "-auto-approve"}, args...))
	}

	return terraformWithProgress("refresh", newArgs)
}
