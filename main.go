package main

import (
	"os"

	"github.com/rethil/fast-nav/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
