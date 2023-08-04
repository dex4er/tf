package progress

import (
	"fmt"

	"github.com/mitchellh/colorstring"
)

var operation2symbol = map[string]string{"R": "=", "C": "+", "M": "~", "D": "-"}
var operation2color = map[string]string{"R": "cyan", "C": "green", "M": "yellow", "D": "red"}

func refreshDots(line string, resource string, operation string) {
	if NoColor {
		fmt.Print("^")
	} else {
		colorstring.Print("[blue]^")
	}
}

func startDots(line string, resource string, operation string) {
	if NoColor {
		fmt.Print(".")
	} else {
		colorstring.Print("[" + operation2color[operation] + "].")
	}
}

func stillDots(line string, resource string, operation string) {
	startDots(line, resource, operation)
}

func stopDots(line string, resource string, operation string) {
	if NoColor {
		fmt.Print(operation2symbol[operation])
	} else {
		colorstring.Print("[" + operation2color[operation] + "]" + operation2symbol[operation])
	}
}
