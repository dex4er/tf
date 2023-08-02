package run

import "fmt"

func Help(args []string) error {
	err := Terraform("", args)

	fmt.Println()
	fmt.Println("Extra commands:")
	fmt.Println("  upgrade       The same as terraform init -upgrade")

	return err
}
