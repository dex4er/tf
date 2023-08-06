package util

import (
	"fmt"
	"os"
	"time"

	"github.com/awoodbeck/strftime"
)

// Opens log file if `TF_OUTPUT_PATH` environment variable is set
func OpenOutputFile() (*os.File, error) {
	if outputPath := os.Getenv("TF_OUTPUT_PATH"); outputPath != "" {
		now := time.Now()
		filename := strftime.Format(&now, outputPath)
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("opening the log file: %w", err)
		}
		return file, nil
	}

	return nil, nil
}
