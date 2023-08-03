package progress

func Refresh(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counter":
		refreshCounter(line, resource, operation)
	}
}

func Start(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counter":
		startCounter(line, resource, operation)
	}
}

func Stop(progressFormat string, line string, resource string, operation string) {
	switch progressFormat {
	case "counter":
		stopCounter(line, resource, operation)
	}
}
