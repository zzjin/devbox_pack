// Package main is the entry point for the DevBox Pack execution plan generator CLI tool.
package main

import (
	"fmt"
	"os"

	"github.com/labring/devbox-pack/pkg/cli"
)

func main() {
	cliHandler := cli.NewCLIApp()
	if err := cliHandler.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
