package main

import (
	"fmt"
	"os"

	"github.com/xyproto/projectinfo"
)

func OutputChunks(dir string) error {
	const printWarnings = true
	pInfo, err := projectinfo.New(dir, printWarnings)
	if err != nil {
		return err
	}
	chunks, err := pInfo.Chunk(true, true)
	if err != nil {
		return err
	}
	for _, chunk := range chunks {
		fmt.Println(chunk)
	}
	return nil
}

func main() {
	// Check for command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: info [directory]")
		os.Exit(1)
	}

	// The first argument should be the directory to scan
	dir := os.Args[1]

	if err := OutputChunks(dir); err != nil {
		fmt.Printf("Failed to output project chunks: %v\n", err)
		os.Exit(1)
	}

}
