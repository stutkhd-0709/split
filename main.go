package main

import (
	"os"

	"github.com/stutkhd-0709/split/cli"
)

func main() {
	exitCode := cli.Main(os.Args)
	os.Exit(exitCode)
}
