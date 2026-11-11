package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// run holds all the startup logic and returns an error.
// Later it will also handle graceful shutdown.
func run() error {
	fmt.Println("taskboard starting")
	return nil
}
