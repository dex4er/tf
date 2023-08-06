package run

import "fmt"

func Version(args []string, version string) error {
	fmt.Println("tf", "v"+version)

	return Terraform("version", args)
}
