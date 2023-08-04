package console

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/term"
)

const defaultCols = 80

var NoColor = false
var Cols = getCols()

func Print(msg string) {
	if strings.HasSuffix(msg, "\n") {
		msg = strings.TrimSuffix(msg, "\n")
		lenMsg := lenWithoutColors(msg)
		if lenMsg < Cols {
			fmt.Println(msg, strings.Repeat(" ", Cols-1-lenMsg))
		} else {
			fmt.Println(msg)
		}
	} else if strings.HasSuffix(msg, "\r") {
		msg = strings.TrimSuffix(msg, "\r")
		lenMsg := lenWithoutColors(msg)
		if lenMsg < Cols {
			fmt.Print(msg, strings.Repeat(" ", Cols-1-lenMsg), "\r")
		} else {
			fmt.Print(msg, "\r")
		}
	} else {
		fmt.Print(msg)
	}
}

func getCols() int {
	width, _, err := term.GetSize(0)
	if err != nil {
		return defaultCols
	}
	return width
}

func lenWithoutColors(msg string) int {
	re := regexp.MustCompile(`\033\[\d+m`)
	re.ReplaceAllString(msg, "")
	return len(msg)
}
