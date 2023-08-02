package run

import "fmt"

func Version(args []string, version string) error {
	fmt.Println("tf", version)

	return Terraform("version", args)
}
