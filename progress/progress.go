package progress

import (
	"github.com/dex4er/tf/progress/counters"
	"github.com/dex4er/tf/progress/dots"
	"github.com/dex4er/tf/progress/fan"
	"github.com/dex4er/tf/progress/verbose"
)

const (
	Counters = "counters"
	Dots     = "dots"
	Fan      = "fan"
	Quiet    = "quiet"
	Verbatim = "verbatim"
	Verbose  = "verbose"
)

// Handles progress indicator when refreshing resources.
func Refresh(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case Counters:
		counters.Refresh(line, resource, operation)
	case Dots:
		dots.Refresh(line, resource, operation)
	case Fan:
		fan.Refresh(line, resource, operation)
	case Verbose:
		verbose.Refresh(line, resource, operation)
	}
}

// Handles progress indicator when preparing import of resources.
func PreparingImport(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case Counters:
		counters.PreparingImport(line, resource, operation)
	case Dots:
		dots.PreparingImport(line, resource, operation)
	case Fan:
		fan.PreparingImport(line, resource, operation)
	case Verbose:
		verbose.PreparingImport(line, resource, operation)
	}
}

// Handles progress indicator when operation on resource starts.
func Start(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case Counters:
		counters.Start(line, resource, operation)
	case Dots:
		dots.Start(line, resource, operation)
	case Fan:
		fan.Start(line, resource, operation)
	case Verbose:
		verbose.Start(line, resource, operation)
	}
}

// Handles progress indicator when operation on resource still is in
func Still(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case Counters:
		counters.Still(line, resource, operation)
	case Dots:
		dots.Still(line, resource, operation)
	case Fan:
		fan.Still(line, resource, operation)
	case Verbose:
		verbose.Still(line, resource, operation)
	}
}

// Handles progress indicator when operation on resource ends.
func Stop(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case Counters:
		counters.Stop(line, resource, operation)
	case Dots:
		dots.Stop(line, resource, operation)
	case Fan:
		fan.Stop(line, resource, operation)
	case Verbose:
		verbose.Stop(line, resource, operation)
	}
}
