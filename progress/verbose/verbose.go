package verbose

import (
	"fmt"

	"github.com/mitchellh/colorstring"

	"github.com/dex4er/tf/console"
)

var refreshing = 0
var started = map[string]int{"Read": 0, "Import": 0, "Creat": 0, "Destr": 0, "Modif": 0}
var stopped = map[string]int{"Read": 0, "Import": 0, "Creat": 0, "Destr": 0, "Modif": 0}

func Refreshing(line string, resource string, operation string) {
	refreshing += 1
	show(line)
}

func PreparingImport(line string, resource string, operation string) {
	refreshing += 1
	show(line)
}

func Start(line string, resource string, operation string) {
	started[operation] += 1
	show(line)
}

func Still(line string, resource string, operation string) {
	show(line)
}

func Stop(line string, resource string, operation string) {
	stopped[operation] += 1
	show(line)
}

func show(line string) {
	s := fmt.Sprintf("^%d", refreshing)
	r := fmt.Sprintf("=%d/%d", stopped["Read"], started["Read"])
	i := fmt.Sprintf("&%d/%d", stopped["Import"], started["Import"])
	c := fmt.Sprintf("+%d/%d", stopped["Creat"], started["Creat"])
	m := fmt.Sprintf("~%d/%d", stopped["Modif"], started["Modif"])
	d := fmt.Sprintf("-%d/%d", stopped["Destr"], started["Destr"])

	if console.NoColor {
		console.Printf("%s %s %s %s %s %s %s\n", s, r, i, c, m, d, line)
	} else {
		console.Printf(colorstring.Color("[blue]%s[reset] [cyan]%s[reset] [dark_gray]%s[reset] [green]%s[reset] [yellow]%s[reset] [red]%s[reset] %s")+"\n", s, r, i, c, m, d, line)
	}
}
