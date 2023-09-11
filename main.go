package main

import (
	"os"

	"github.com/stutkhd-0709/split/cli"
)

func main() {
	cli := &cli.CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	// os.Args[1:]でいいのか確認
	os.Exit(cli.RunCommand(os.Args[1:]))
}
