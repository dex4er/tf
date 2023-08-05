package run

import "fmt"

func Help(args []string) error {
	// should provoke an error
	err := Terraform("", args)

	fmt.Println()
	fmt.Println("Extra commands:")
	fmt.Println("  mv            The same as terraform state mv")
	fmt.Println("  rm            The same as terraform state rm")
	fmt.Println("  upgrade       The same as terraform init -upgrade")

	return err
}
