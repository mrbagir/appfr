package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	switch os.Args[1] {
	case "init":
		initCmd()
	}
}

// create project structure and files from template.
func initCmd() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: init <module-name>")
		return
	}

	fmt.Println("Coming Soon")
}
