package progress

func Refresh(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counters":
		refreshCounters(line, resource, operation)
	case "dots":
		refreshDots(line, resource, operation)
	}
}

func Start(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counters":
		startCounters(line, resource, operation)
	case "dots":
		startDots(line, resource, operation)
	}
}

func Still(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counters":
		stillCounters(line, resource, operation)
	case "dots":
		stillDots(line, resource, operation)
	}
}

func Stop(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counters":
		stopCounters(line, resource, operation)
	case "dots":
		stopDots(line, resource, operation)
	}
}
