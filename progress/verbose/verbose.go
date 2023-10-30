package verbose

import (
	"fmt"

	"github.com/mitchellh/colorstring"

	"github.com/dex4er/tf/console"
)

var refreshing = 0
var started = map[string]int{"R": 0, "I": 0, "C": 0, "D": 0, "M": 0}
var stopped = map[string]int{"R": 0, "I": 0, "C": 0, "D": 0, "M": 0}

func Refreshing(line string, resource string, operation string) {
	refreshing += 1
	show(line, resource, operation)
}

func PreparingImport(line string, resource string, operation string) {
	refreshing += 1
	show(line, resource, operation)
}

func Start(line string, resource string, operation string) {
	started[operation] += 1
	show(line, resource, operation)
}

func Still(line string, resource string, operation string) {
	show(line, resource, operation)
}

func Stop(line string, resource string, operation string) {
	stopped[operation] += 1
	show(line, resource, operation)
}

func show(line string, resource string, operation string) {
	s := fmt.Sprintf("^%d", refreshing)
	r := fmt.Sprintf("=%d/%d", stopped["R"], started["R"])
	i := fmt.Sprintf("&%d/%d", stopped["I"], started["I"])
	c := fmt.Sprintf("+%d/%d", stopped["C"], started["C"])
	m := fmt.Sprintf("~%d/%d", stopped["M"], started["M"])
	d := fmt.Sprintf("-%d/%d", stopped["D"], started["D"])

	if console.NoColor {
		console.Printf("%s %s %s %s %s %s %s\n", s, r, i, c, m, d, line)
	} else {
		console.Printf(colorstring.Color("[blue]%s[reset] [cyan]%s[reset] [dark_gray]%s[reset] [green]%s[reset] [yellow]%s[reset] [red]%s[reset] %s")+"\n", s, r, i, c, m, d, line)
	}
}
