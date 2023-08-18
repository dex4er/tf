package run

import (
	"strings"

	"github.com/dex4er/tf/util"
)

func Destroy(args []string) error {
	newArgs := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			if util.ReplaceFirstTwoDashes(arg) == "-target" {
				continue
			}
			newArgs = append(newArgs, arg)
		} else {
			newArgs = append(newArgs, "-target="+util.AddQuotes(arg))
		}
	}

	return terraformWithProgress("destroy", newArgs)
}
