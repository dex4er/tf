package progress

import (
	"github.com/dex4er/tf/progress/counters"
	"github.com/dex4er/tf/progress/dots"
	"github.com/dex4er/tf/progress/fan"
	"github.com/dex4er/tf/progress/verbose"
)

// Handles progress indicator when refreshing resources.
func Refresh(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counters":
		counters.Refresh(line, resource, operation)
	case "dots":
		dots.Refresh(line, resource, operation)
	case "fan":
		fan.Refresh(line, resource, operation)
	case "verbose":
		verbose.Refresh(line, resource, operation)
	}
}

// Handles progress indicator when operation on resource starts.
func Start(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counters":
		counters.Start(line, resource, operation)
	case "dots":
		dots.Start(line, resource, operation)
	case "fan":
		fan.Start(line, resource, operation)
	case "verbose":
		verbose.Start(line, resource, operation)
	}
}

// Handles progress indicator when operation on resource still is in progress.
func Still(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counters":
		counters.Still(line, resource, operation)
	case "dots":
		dots.Still(line, resource, operation)
	case "fan":
		fan.Still(line, resource, operation)
	case "verbose":
		verbose.Still(line, resource, operation)
	}
}

// Handles progress indicator when operation on resource ends.
func Stop(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counters":
		counters.Stop(line, resource, operation)
	case "dots":
		dots.Stop(line, resource, operation)
	case "fan":
		fan.Stop(line, resource, operation)
	case "verbose":
		verbose.Stop(line, resource, operation)
	}
}
