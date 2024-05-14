# projectinfo

[![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/projectinfo)](https://goreportcard.com/report/github.com/xyproto/projectinfo) [![GoDoc](https://godoc.org/github.com/xyproto/projectinfo?status.svg)](https://godoc.org/github.com/xyproto/projectinfo) [![License](https://img.shields.io/badge/license-BSD-green.svg?style=flat)](https://raw.githubusercontent.com/xyproto/projectinfo/main/LICENSE)

Given a directory of source code, gather all sorts of info and output is as chunks of JSON.

Example usage:

```go
package main

import (
    "fmt"
    "os"

    "github.com/xyproto/projectinfo"
)

func OutputChunks(dir string) error {
    pi, err := projectinfo.New(dir)
    if err != nil {
        return err
    }
    chunks, err := pi.Chunk(true, true)
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
```

## General info

* Version: 1.1.1
* License: BSD-3
