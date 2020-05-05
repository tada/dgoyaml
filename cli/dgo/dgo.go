// Package main contains the main dgo CLI command.
package main

import (
	"os"

	"github.com/tada/dgoyaml/cli"
)

func main() {
	// Could use spf13.cobra here but it brings in a fairly large and undesired set of transitive dependencies
	cli.Dgo(os.Stdout, os.Stderr).Do(os.Args[1:])
}
