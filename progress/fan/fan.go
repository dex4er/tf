package fan

import (
	"fmt"

	"github.com/mitchellh/colorstring"

	"github.com/dex4er/tf/console"
	"github.com/dex4er/tf/progress/operations"
)

var fanTicks = []string{"-", `\`, "|", "/"}
var fanIndex = 0

var pending = map[string]string{}

func Refresh(line string, resource string, operation string) {
	show()
}

func PreparingImport(line string, resource string, operation string) {
	show()
}

func Start(line string, resource string, operation string) {
	pending[resource] = operation
	show()
}

func Still(line string, resource string, operation string) {
	pending[resource] = operation
	show()
}

func Stop(line string, resource string, operation string) {
	delete(pending, resource)
	show()
}

func show() {
	fanIndex = (fanIndex + 1) % len(fanTicks)
	f := fanTicks[fanIndex]
	i := 0
	for _, v := range pending {
		if console.NoColor {
			f += operations.Operation2symbol[v]
		} else {
			f += colorstring.Color(fmt.Sprintf("[%s]%s", operations.Operation2color[v], operations.Operation2symbol[v]))
		}
		if i > console.Cols-2 {
			break
		}
		i++
	}
	console.Print(f + "\r")
}
