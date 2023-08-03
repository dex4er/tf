package console

import (
	"fmt"
	"strings"

	"golang.org/x/term"
)

const defaultCols = 80

var cols = getCols()

func Print(msg string) {
	if strings.HasSuffix(msg, "\n") {
		msg = strings.TrimSuffix(msg, "\n")
		lenMsg := len(msg)
		if lenMsg < cols {
			fmt.Print(msg)
			fmt.Println(strings.Repeat(" ", cols-1-lenMsg))
		} else {
			fmt.Println(msg)
		}
	} else {
		fmt.Print(msg)
	}
}

func getCols() int {
	width, _, err := term.GetSize(0)
	if err != nil {
		return defaultCols
	}
	return width
}
