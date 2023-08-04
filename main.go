package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dex4er/tf/run"
)

var version = "dev"

var emptyStringsList = make([]string, 0)

func main() {
	var err error

	if len(os.Args) < 2 {
		err = run.Help(emptyStringsList)
	} else {
		command := os.Args[1]
		args := os.Args[2:]

		switch command {
		case "apply":
			err = run.Apply(args)
		case "destroy":
			err = run.Destroy(args)
		case "init":
			err = run.Init(args)
		case "list":
			err = run.List(args)
		case "plan":
			err = run.Plan(args)
		case "upgrade":
			err = run.Upgrade(args)
		case "version":
			err = run.Version(args, version)
		default:
			err = run.Terraform(command, args)
		}
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			os.Exit(exitCode)
		} else {
			fmt.Println("Error:", err)
			os.Exit(2)
		}
	}
}
