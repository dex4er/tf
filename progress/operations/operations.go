package operations

const (
	Refreshing = "refresh"
	Reading    = "read"
	Opening    = "open"
	Closing    = "clos"
	Importing  = "import"
	Creating   = "creat"
	Modifying  = "modif"
	Destroying = "destr"
)

var Operation2symbol = map[string]string{Refreshing: "^", Reading: "=", Opening: "<", Closing: ">", Importing: "&", Creating: "+", Modifying: "~", Destroying: "-"}
var Operation2color = map[string]string{Refreshing: "blue", Reading: "cyan", Opening: "blue", Closing: "blue", Importing: "dark_gray", Creating: "green", Modifying: "yellow", Destroying: "red"}
