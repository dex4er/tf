package run

import (
	"os"
	"strings"

	"github.com/dex4er/tf/util"
)

func Apply(args []string) error {
	newArgs := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if util.ReplaceFirstTwoDashes(arg) == "-target" {
				continue
			}
			newArgs = append(newArgs, arg)
		} else {
			if _, err := os.Stat(arg); err == nil {
				newArgs = append(newArgs, arg)
			} else {
				newArgs = append(newArgs, "-target="+util.AddQuotes(arg))
			}
		}
	}

	return terraformWithProgress("apply", newArgs)
}
