package main

import (
	"os"
	"testing"
)

// The sole purpose of this test it to get 100% code coverage
func Test_main(_ *testing.T) {
	os.Args = []string{`dgo`, `--help`}
	main()
}
