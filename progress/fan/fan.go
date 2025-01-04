package fan

import (
	"fmt"

	"github.com/dex4er/tf/console"
	"github.com/mitchellh/colorstring"
)

var fanTicks = []string{"-", `\`, "|", "/"}
var fanIndex = 0

var operation2symbol = map[string]string{"Read": "=", "Import": "&", "Creat": "+", "Modif": "~", "Destr": "-"}
var operation2color = map[string]string{"Read": "cyan", "Import": "dark_gray", "Creat": "green", "Modif": "yellow", "Destr": "red"}

var operations = map[string]string{}

func Refreshing(line string, resource string, operation string) {
	show()
}

func PreparingImport(line string, resource string, operation string) {
	show()
}

func Start(line string, resource string, operation string) {
	operations[resource] = operation
	show()
}

func Still(line string, resource string, operation string) {
	operations[resource] = operation
	show()
}

func Stop(line string, resource string, operation string) {
	delete(operations, resource)
	show()
}

func show() {
	fanIndex = (fanIndex + 1) % len(fanTicks)
	f := fanTicks[fanIndex]
	i := 0
	for _, v := range operations {
		if console.NoColor {
			f += operation2symbol[v]
		} else {
			f += colorstring.Color(fmt.Sprintf("[%s]%s", operation2color[v], operation2symbol[v]))
		}
		if i > console.Cols-2 {
			break
		}
		i++
	}
	console.Print(f + "\r")
}
