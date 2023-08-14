package main

import (
	"os"

	"github.com/stutkhd-0709/enable_bootcamp/cli"
)

func main() {
	exitCode := cli.Main()
	os.Exit(exitCode)
}
