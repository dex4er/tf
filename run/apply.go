package run

import (
	"strings"

	"github.com/dex4er/tf/util"
)

func Apply(args []string) error {
	newArgs := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			newArgs = append(newArgs, arg)
		} else {
			newArgs = append(newArgs, "-target="+util.AddQuotes(arg))
		}
	}

	return terraformWithProgress("apply", newArgs)
}
