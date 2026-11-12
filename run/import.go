package run

import (
	"strings"

	"github.com/dex4er/tf/util"
)

func Import(args []string) error {
	newArgs, resource, noShow := importArgs(args)

	if err := terraformWithProgress("import", newArgs); err != nil {
		return err
	}

	if !noShow {
		return Show([]string{resource})
	}

	return nil
}

// The first positional argument is the address and the rest is the ID, joined
// with spaces.
func importArgs(args []string) (newArgs []string, resource string, noShow bool) {
	newArgs = []string{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-no-show" {
			noShow = true
			continue
		}
		if strings.HasPrefix(arg, "-") {
			optionArgs, n := optionArgs(args, i)
			newArgs = append(newArgs, optionArgs...)
			i += n - 1
			continue
		}
		resource = util.AddQuotes(arg)
		newArgs = append(newArgs, resource)
		if len(args) > i+1 {
			newArgs = append(newArgs, strings.Join(args[i+1:], " "))
		}
		break
	}

	return newArgs, resource, noShow
}
