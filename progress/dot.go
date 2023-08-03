package progress

import (
	"fmt"
)

var operation2symbol = map[string]string{"R": "=", "C": "+", "M": "~", "D": "-"}

func refreshDot(line string, resource string, operation string) {
	fmt.Print("^")
}

func startDot(line string, resource string, operation string) {
	fmt.Print(".")
}

func stillDot(line string, resource string, operation string) {
	fmt.Print(".")
}

func stopDot(line string, resource string, operation string) {
	fmt.Print(operation2symbol[operation])
}
