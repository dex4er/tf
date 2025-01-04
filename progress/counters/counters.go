package counters

import (
	"fmt"

	"github.com/mitchellh/colorstring"

	"github.com/dex4er/tf/console"
	"github.com/dex4er/tf/progress/operations"
)

var Refreshing = 0
var Started = map[string]int{operations.Reading: 0, operations.Opening: 0, operations.Closing: 0, operations.Importing: 0, operations.Creating: 0, operations.Destroying: 0}
var Stopped = map[string]int{operations.Reading: 0, operations.Opening: 0, operations.Closing: 0, operations.Importing: 0, operations.Creating: 0, operations.Destroying: 0}

func Refresh(line string, resource string, operation string) {
	Refreshing += 1
	show(line)
}

func PreparingImport(line string, resource string, operation string) {
	Refreshing += 1
	show(line)
}

func Start(line string, resource string, operation string) {
	Started[operation] += 1
	show(line)
}

func Still(line string, resource string, operation string) {
	show(line)
}

func Stop(line string, resource string, operation string) {
	Stopped[operation] += 1
	show(line)
}

func show(line string) {
	counters, countersLength := Counters()

	maxLine := max(console.Cols-countersLength, 0)
	lineTrimmed := line[:min(len(line), maxLine)]

	console.Printf("%s%s\r", counters, lineTrimmed)
}

func Counters() (string, int) {
	countersLength := 0
	counters := ""

	if c, l := counter(operations.Refreshing); l > 0 {
		countersLength += l
		counters += c
	}
	if c, l := counter(operations.Reading); l > 0 {
		countersLength += l
		counters += c
	}
	if c, l := counter(operations.Opening); l > 0 {
		countersLength += l
		counters += c
	}
	if c, l := counter(operations.Closing); l > 0 {
		countersLength += l
		counters += c
	}
	if c, l := counter(operations.Importing); l > 0 {
		countersLength += l
		counters += c
	}
	if c, l := counter(operations.Creating); l > 0 {
		countersLength += l
		counters += c
	}
	if c, l := counter(operations.Destroying); l > 0 {
		countersLength += l
		counters += c
	}
	if c, l := counter(operations.Modifying); l > 0 {
		countersLength += l
		counters += c
	}

	return counters, countersLength
}

func counter(operation string) (string, int) {
	c := ""
	if operation == operations.Refreshing {
		if Refreshing > 0 {
			c = fmt.Sprintf("%s%d ", operations.Operation2symbol[operation], Refreshing)
		}
	} else {
		if Stopped[operation]+Started[operation] > 0 {
			c = fmt.Sprintf("%s%d/%d ", operations.Operation2symbol[operation], Stopped[operation], Started[operation])
		}
	}
	l := len(c)
	if l > 0 && !console.NoColor {
		return colorstring.Color("[" + operations.Operation2color[operation] + "]" + c + "[reset]"), l
	}
	return c, l
}
