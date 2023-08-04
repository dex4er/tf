package progress

import (
	"fmt"
)

var operation2symbol = map[string]string{"R": "=", "C": "+", "M": "~", "D": "-"}

func refreshDots(line string, resource string, operation string) {
	fmt.Print("^")
}

func startDots(line string, resource string, operation string) {
	fmt.Print(".")
}

func stillDots(line string, resource string, operation string) {
	fmt.Print(".")
}

func stopDots(line string, resource string, operation string) {
	fmt.Print(operation2symbol[operation])
}
