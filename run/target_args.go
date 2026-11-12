package run

import (
	"os"
	"strings"

	"github.com/dex4er/tf/util"
)

// files keeps an argument that names an existing file positional, as a plan
// file, instead of taking it for an address.
func targetArgs(args []string, files bool) []string {
	newArgs := []string{}

	for _, arg := range args {
		if address, ok := util.TargetAddress(arg); ok {
			newArgs = append(newArgs, "-target="+util.AddQuotes(address))
			continue
		}
		if strings.HasPrefix(arg, "-") {
			// the address of `-target ADDRESS` follows as a positional argument
			if util.ReplaceFirstTwoDashes(arg) == "-target" {
				continue
			}
			newArgs = append(newArgs, arg)
			continue
		}
		if files {
			if _, err := os.Stat(arg); err == nil {
				newArgs = append(newArgs, arg)
				continue
			}
		}
		newArgs = append(newArgs, "-target="+util.AddQuotes(arg))
	}

	return newArgs
}

// The plan file of `-plan=FILE` goes last, as terraform reads no option after
// it.
func applyArgs(args []string) []string {
	otherArgs := []string{}
	planFile := ""

	for _, arg := range args {
		if file, ok := strings.CutPrefix(util.ReplaceFirstTwoDashes(arg), "-plan="); ok {
			planFile = file
			continue
		}
		otherArgs = append(otherArgs, arg)
	}

	newArgs := targetArgs(otherArgs, true)

	if planFile != "" {
		newArgs = append(newArgs, planFile)
	}

	return newArgs
}
