package run

import "fmt"

func Version(args []string, version string) error {
	err := Terraform("version", args)

	fmt.Println("+ tf", "v"+version)

	return err
}
