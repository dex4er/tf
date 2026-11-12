package run

import (
	"os"

	"github.com/dex4er/tf/progress"
	"github.com/dex4er/tf/util"
)

type progressOptions struct {
	planFormat     string
	progressFormat string
	noColor        bool
	noOutputs      bool
	// the arguments for terraform, without the options of tf
	args []string
}

// The environment is read at each call rather than once, so that a test can set
// it.
func parseProgressOptions(args []string) progressOptions {
	opts := progressOptions{
		planFormat:     Short,
		progressFormat: progress.Counters,
		noOutputs:      true,
		args:           []string{},
	}

	if os.Getenv("TF_IN_AUTOMATION") != "" {
		opts.progressFormat = progress.Verbose
	}

	if util.IsAgent() {
		opts.progressFormat = progress.Quiet
	}

	if format := os.Getenv("TF_PLAN_FORMAT"); format != "" {
		opts.planFormat = format
	}

	if format := os.Getenv("TF_PROGRESS_FORMAT"); format != "" {
		opts.progressFormat = format
	}

	for _, arg := range args {
		switch util.ReplaceFirstTwoDashes(arg) {
		case "-compact":
			opts.planFormat = Compact
		case "-counter", "-counters":
			opts.progressFormat = progress.Counters
		case "-dot", "-dots":
			opts.progressFormat = progress.Dots
		case "-fan":
			opts.progressFormat = progress.Fan
		case "-full":
			opts.planFormat = Full
		case "-no-color":
			opts.noColor = true
			opts.args = append(opts.args, arg)
		case "-no-colors":
			opts.noColor = true
			opts.args = append(opts.args, "-no-color")
		case "-no-output", "-no-outputs":
			opts.noOutputs = true
		case "-no-output=false", "-no-outputs=false":
			opts.noOutputs = false
		case "-quiet":
			opts.progressFormat = progress.Quiet
		case "-short":
			opts.planFormat = Short
		case "-verbatim":
			opts.planFormat = Full
			opts.progressFormat = progress.Verbatim
		case "-verbose":
			opts.progressFormat = progress.Verbose
		default:
			opts.args = append(opts.args, arg)
		}
	}

	// first, as terraform reads no option after a positional argument, as the
	// plan file of apply
	if !opts.noColor && util.NoColor() {
		opts.noColor = true
		opts.args = append([]string{"-no-color"}, opts.args...)
	}

	return opts
}
