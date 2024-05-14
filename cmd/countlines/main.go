package main

import (
	"fmt"
	"os"

	"github.com/xyproto/projectinfo"
)

func CountLines(filename string) (int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	return projectinfo.CountLines(string(data))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: countlines [filename]")
		os.Exit(1)
	}
	filename := os.Args[1]

	count, err := CountLines(filename)
	if err != nil {
		fmt.Printf("Failed to count lines in %s: %v\n", filename, err)
		os.Exit(1)
	}

	fmt.Printf("%d\n", count)
}
