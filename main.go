package main

import (
	"os"

	"github.com/stutkhd-0709/split/cli"
)

const (
	ExitOK int = 0
	ExitNG int = 1
)

func main() {
	cli := &cli.CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	os.Exit(cli.RunCommand(os.Args))
}
