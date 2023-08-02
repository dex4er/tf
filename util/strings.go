package util

import (
	"regexp"
	"strings"
)

func AddQuotes(input string) string {
	re := regexp.MustCompile(`\[([a-z_][^\]]*)\]`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		return `["` + re.FindStringSubmatch(match)[1] + `"]`
	})
}

func IsEmptyLine(line string) bool {
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	line = strings.TrimSuffix(line, ColorReset)
	return line == ""
}

func StartsWith(s string, prefix rune) bool {
	if len(s) == 0 {
		return false
	}
	return []rune(s)[0] == prefix
}
