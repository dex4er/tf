package util

import (
	"regexp"
)

func AddQuotes(input string) string {
	re := regexp.MustCompile(`\[([A-Za-z_][^\]]*)\]`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		return `["` + re.FindStringSubmatch(match)[1] + `"]`
	})
}

func IsEmptyLine(line string) bool {
	re := regexp.MustCompile(`^(\033\[\d+m[╷╵]?)*\r?\n?$`)
	return re.MatchString(line)
}

func RemoveColors(line string) string {
	re := regexp.MustCompile(`\033\[\d+m`)
	return re.ReplaceAllString(line, "")
}

func StartsWith(s string, prefix rune) bool {
	if len(s) == 0 {
		return false
	}
	return []rune(s)[0] == prefix
}
