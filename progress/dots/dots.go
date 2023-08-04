package dots

import (
	"fmt"

	"github.com/mitchellh/colorstring"

	"github.com/dex4er/tf/console"
)

var operation2symbol = map[string]string{"R": "=", "C": "+", "M": "~", "D": "-"}
var operation2color = map[string]string{"R": "cyan", "C": "green", "M": "yellow", "D": "red"}

func Refresh(line string, resource string, operation string) {
	if console.NoColor {
		fmt.Print("^")
	} else {
		colorstring.Print("[blue]^")
	}
}

func Start(line string, resource string, operation string) {
	if console.NoColor {
		fmt.Print(".")
	} else {
		colorstring.Print("[" + operation2color[operation] + "].")
	}
}

func Still(line string, resource string, operation string) {
	Start(line, resource, operation)
}

func Stop(line string, resource string, operation string) {
	if console.NoColor {
		fmt.Print(operation2symbol[operation])
	} else {
		colorstring.Print("[" + operation2color[operation] + "]" + operation2symbol[operation])
	}
}
