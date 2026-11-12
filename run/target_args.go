package run

import (
	"os"
	"strings"

	"github.com/dex4er/tf/util"
)

// The options of terraform and tofu that take a value, which can follow as the
// next argument instead of after `=`, and which is then no address. Their
// boolean options take no value that way.
var valueOptions = map[string]bool{
	"-backup":              true,
	"-config":              true,
	"-deprecation":         true,
	"-exclude":             true,
	"-exclude-file":        true,
	"-generate-config-out": true,
	"-invoke":              true,
	"-json-into":           true,
	"-lock-timeout":        true,
	"-out":                 true,
	"-parallelism":         true,
	"-plan":                true,
	"-policies":            true,
	"-provider":            true,
	"-replace":             true,
	"-state":               true,
	"-state-out":           true,
	"-target-file":         true,
	"-var":                 true,
	"-var-file":            true,
}

// The options besides `-target` whose value is an address, which gets its
// quotes back.
var addressOptions = map[string]bool{
	"-exclude": true,
	"-invoke":  true,
	"-replace": true,
}

// Returns the arguments of an option at args[i]: the option, with the address
// of an address option quoted, and its value when it follows as the next
// argument. n is how many arguments it takes.
func optionArgs(args []string, i int) (optionArgs []string, n int) {
	option := util.ReplaceFirstTwoDashes(args[i])

	if name, value, ok := strings.Cut(option, "="); ok {
		if addressOptions[name] {
			return []string{name + "=" + util.AddQuotes(value)}, 1
		}
		return []string{args[i]}, 1
	}

	if valueOptions[option] && i+1 < len(args) {
		value := args[i+1]
		if addressOptions[option] {
			value = util.AddQuotes(value)
		}
		return []string{args[i], value}, 2
	}

	return []string{args[i]}, 1
}

// files keeps an argument that names an existing file positional, as a plan
// file, instead of taking it for an address.
func targetArgs(args []string, files bool) []string {
	newArgs := []string{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if address, ok := util.TargetAddress(arg); ok {
			newArgs = append(newArgs, "-target="+util.AddQuotes(address))
			continue
		}
		if strings.HasPrefix(arg, "-") {
			// the address of `-target ADDRESS` follows as a positional argument
			if util.ReplaceFirstTwoDashes(arg) == "-target" {
				continue
			}
			optionArgs, n := optionArgs(args, i)
			newArgs = append(newArgs, optionArgs...)
			i += n - 1
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

// The plan file of `-plan=FILE` or `-plan FILE` goes last, as terraform reads
// no option after it.
func applyArgs(args []string) []string {
	otherArgs := []string{}
	planFile := ""

	for i := 0; i < len(args); i++ {
		option := util.ReplaceFirstTwoDashes(args[i])
		if file, ok := strings.CutPrefix(option, "-plan="); ok {
			planFile = file
			continue
		}
		if option == "-plan" && i+1 < len(args) {
			planFile = args[i+1]
			i++
			continue
		}
		otherArgs = append(otherArgs, args[i])
	}

	newArgs := targetArgs(otherArgs, true)

	if planFile != "" {
		newArgs = append(newArgs, planFile)
	}

	return newArgs
}

// Splits the arguments of the commands that take addresses without `-target`
// into the options, with their values, and the addresses, quoted.
func splitAddresses(args []string) (options []string, addresses []string) {
	options = []string{}
	addresses = []string{}

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			optionArgs, n := optionArgs(args, i)
			options = append(options, optionArgs...)
			i += n - 1
			continue
		}
		addresses = append(addresses, util.AddQuotes(args[i]))
	}

	return options, addresses
}
