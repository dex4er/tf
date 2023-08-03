package progress

func Refresh(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counter":
		refreshCounter(line, resource, operation)
	case "dot":
		refreshDot(line, resource, operation)
	}
}

func Start(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counter":
		startCounter(line, resource, operation)
	case "dot":
		startDot(line, resource, operation)
	}
}

func Still(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counter":
		stillCounter(line, resource, operation)
	case "dot":
		stillDot(line, resource, operation)
	}
}

func Stop(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counter":
		stopCounter(line, resource, operation)
	case "dot":
		stopDot(line, resource, operation)
	}
}
