package util

import (
	"fmt"
	"os"
)

func OpenLogfile() (*os.File, error) {
	if logFilename := os.Getenv("TF_LOG_FILE"); logFilename != "" {
		file, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("opening the log file: %w", err)
		}
		return file, nil
	}

	return nil, nil
}
