package progress

import (
	"fmt"

	"github.com/mitchellh/colorstring"

	"github.com/dex4er/tf/console"
	"github.com/dex4er/tf/util"
)

var refreshed = 0
var started = map[string]int{"R": 0, "C": 0, "D": 0, "M": 0}
var stopped = map[string]int{"R": 0, "C": 0, "D": 0, "M": 0}

func refreshCounters(line string, resource string, operation string) {
	refreshed += 1
	showCounters(line, resource, operation)
}

func startCounters(line string, resource string, operation string) {
	started[operation] += 1
	showCounters(line, resource, operation)
}

func stillCounters(line string, resource string, operation string) {
	showCounters(line, resource, operation)
}

func stopCounters(line string, resource string, operation string) {
	stopped[operation] += 1
	showCounters(line, resource, operation)
}

func showCounters(line string, resource string, operation string) {
	s := fmt.Sprintf("^%d", refreshed)
	r := fmt.Sprintf("=%d/%d", stopped["R"], started["R"])
	c := fmt.Sprintf("+%d/%d", stopped["C"], started["C"])
	m := fmt.Sprintf("~%d/%d", stopped["M"], started["M"])
	d := fmt.Sprintf("-%d/%d", stopped["D"], started["D"])

	maxLine := util.MaxInt(console.Cols-len(s)-len(r)-len(c)-len(m)-len(d)-6, 0)
	l := line[:util.MinInt(len(line), maxLine)]

	if NoColor {
		fmt.Printf("%s %s %s %s %s %s\r", s, r, c, m, d, l)
	} else {
		colorstring.Printf("[blue]%s[reset] [cyan]%s[reset] [green]%s[reset] [yellow]%s[reset] [red]%s[reset] %s\r", s, r, c, m, d, l)
	}
}
