package console

import (
	"fmt"
	"strings"

	"golang.org/x/term"
)

const defaultCols = 80

// Flag enabled if `-no-color` option is used.
var NoColor = false

// Number of columns detected from Stdin or 80 by default.
var Cols = getCols()

var previousWasCR = false

// Print string to the console with space padding so progress indicator might be overriden.
func Print(msg string) {
	if previousWasCR {
		fmt.Print(strings.Repeat(" ", Cols-1), "\r")
		previousWasCR = false
	}
	if strings.HasSuffix(msg, "\r") {
		fmt.Print(msg)
		previousWasCR = true
	}
	fmt.Print(msg)
}

// Printf string to the console with space padding so progress indicator might be overriden.
func Printf(format string, a ...interface{}) {
	Print(fmt.Sprintf(format, a...))
}

func getCols() int {
	width, _, err := term.GetSize(0)
	if err != nil {
		return defaultCols
	}
	return width
}
