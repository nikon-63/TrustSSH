package main

import (
	"fmt"
	"os"

	"github.com/nikon-63/TrustSSH/cli/cmd"
)

func main() {
	if err := cmd.Execute(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
