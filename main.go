package main

import (
	"os"

	"github.com/rethil/fn/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}