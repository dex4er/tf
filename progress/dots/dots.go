package dots

import (
	"fmt"

	"github.com/mitchellh/colorstring"

	"github.com/dex4er/tf/console"
	"github.com/dex4er/tf/progress/operations"
)

func Refresh(line string, resource string, operation string) {
	if console.NoColor {
		fmt.Print(operations.Operation2symbol[operation])
	} else {
		console.Print(colorstring.Color("[" + operations.Operation2color[operation] + "]" + operations.Operation2symbol[operation]))
	}
}

func PreparingImport(line string, resource string, operation string) {
	if console.NoColor {
		fmt.Print(operations.Operation2symbol[operation])
	} else {
		console.Print(colorstring.Color("[" + operations.Operation2color[operation] + "]" + operations.Operation2symbol[operation]))
	}
}

func Start(line string, resource string, operation string) {
	if console.NoColor {
		fmt.Print(operations.Operation2symbol[operation])
	} else {
		console.Print(colorstring.Color("[" + operations.Operation2color[operation] + "]."))
	}
}

func Still(line string, resource string, operation string) {
	Start(line, resource, operation)
}

func Stop(line string, resource string, operation string) {
	if console.NoColor {
		fmt.Print(operations.Operation2symbol[operation])
	} else {
		console.Print(colorstring.Color("[" + operations.Operation2color[operation] + "]" + operations.Operation2symbol[operation]))
	}
}
