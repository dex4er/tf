package util

import (
	"strings"
)

func IsEmptyLine(line string) bool {
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	line = strings.TrimSuffix(line, ColorReset)
	return line == ""
}
