package verbose

import (
	"github.com/dex4er/tf/console"
	"github.com/dex4er/tf/progress/counters"
)

func Refresh(line string, resource string, operation string) {
	counters.Refreshing += 1
	show(line)
}

func PreparingImport(line string, resource string, operation string) {
	counters.Refreshing += 1
	show(line)
}

func Start(line string, resource string, operation string) {
	counters.Started[operation] += 1
	show(line)
}

func Still(line string, resource string, operation string) {
	show(line)
}

func Stop(line string, resource string, operation string) {
	counters.Stopped[operation] += 1
	show(line)
}

func show(line string) {
	counters, _ := counters.Counters()

	console.Printf("%s%s\n", counters, line)
}
