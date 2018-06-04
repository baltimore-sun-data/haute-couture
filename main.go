package main

import (
	"fmt"
	"os"

	"github.com/baltimore-sun-data/haute-couture/checker"
)

func main() {
	conf := checker.FromArgs(os.Args[1:])
	if err := conf.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}
}
