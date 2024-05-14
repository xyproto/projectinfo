package main

import (
	"fmt"
	"os"

	"github.com/xyproto/projectinfo"
)

func main() {
	// Check for command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: projectname [directory]")
		os.Exit(1)
	}

	// The first argument should be the directory to scan
	dir := os.Args[1]

	// Use the ReadProjectName function from the projectname package
	name, err := projectinfo.ReadProjectName(dir)
	if err != nil {
		fmt.Printf("Failed to read project name: %v\n", err)
		os.Exit(1)
	}

	// Output the found project name
	fmt.Printf("Project name: %s\n", name)
}
