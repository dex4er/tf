package progress

import (
	"fmt"

	"github.com/dex4er/tf/console"
	"github.com/mitchellh/colorstring"
)

var refreshed = 0
var started = map[string]int{"R": 0, "C": 0, "D": 0, "M": 0}
var stopped = map[string]int{"R": 0, "C": 0, "D": 0, "M": 0}

func refreshCounter(line string, resource string, operation string) {
	refreshed += 1
	showCounter(line, resource, operation)
}

func startCounter(line string, resource string, operation string) {
	started[operation] += 1
	showCounter(line, resource, operation)
}

func stillCounter(line string, resource string, operation string) {
	showCounter(line, resource, operation)
}

func stopCounter(line string, resource string, operation string) {
	stopped[operation] += 1
	showCounter(line, resource, operation)
}

func showCounter(line string, resource string, operation string) {
	s := fmt.Sprintf("^%d", refreshed)
	r := fmt.Sprintf("=%d/%d", stopped["R"], started["R"])
	c := fmt.Sprintf("+%d/%d", stopped["C"], started["C"])
	m := fmt.Sprintf("~%d/%d", stopped["M"], started["M"])
	d := fmt.Sprintf("-%d/%d", stopped["D"], started["D"])

	maxLine := maxInt(console.Cols-len(s)-len(r)-len(c)-len(m)-len(d)-6, 0)
	l := line[:minInt(len(line), maxLine)]

	colorstring.Printf("[blue]%s[reset] [cyan]%s[reset] [green]%s[reset] [yellow]%s[reset] [red]%s[reset] %s\r", s, r, c, m, d, l)
}

func maxInt(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func minInt(x int, y int) int {
	if x < y {
		return x
	}
	return y
}
