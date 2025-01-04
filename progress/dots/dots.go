package dots

import (
	"fmt"

	"github.com/mitchellh/colorstring"

	"github.com/dex4er/tf/console"
)

var operation2symbol = map[string]string{"Read": "=", "Import": "&", "Creat": "+", "Modif": "~", "Destr": "-", "Open": "<", "Clos": ">"}
var operation2color = map[string]string{"Read": "cyan", "Import": "dark_gray", "Creat": "green", "Modif": "yellow", "Destr": "red", "Open": "blue", "Clos": "blue"}

func Refreshing(line string, resource string, operation string) {
	if console.NoColor {
		fmt.Print("^")
	} else {
		console.Print(colorstring.Color("[blue]^"))
	}
}

func PreparingImport(line string, resource string, operation string) {
	if console.NoColor {
		fmt.Print("&")
	} else {
		console.Print(colorstring.Color("[dark_gray]&"))
	}
}

func Start(line string, resource string, operation string) {
	if console.NoColor {
		fmt.Print(".")
	} else {
		console.Print(colorstring.Color("[" + operation2color[operation] + "]."))
	}
}

func Still(line string, resource string, operation string) {
	Start(line, resource, operation)
}

func Stop(line string, resource string, operation string) {
	if console.NoColor {
		console.Print(operation2symbol[operation])
	} else {
		console.Print(colorstring.Color("[" + operation2color[operation] + "]" + operation2symbol[operation]))
	}
}
