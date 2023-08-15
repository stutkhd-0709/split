package cli

import (
	"fmt"
	"os"
)

const (
	ExitOK int = 0
	ExitNG int = 1
)

func Main(args []string) int {
	cli := &CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	if err := cli.RunCommand(os.Args); err != nil {
		fmt.Fprintln(cli.Stderr, "[ERROR] ", err)
		return ExitNG
	}

	return ExitOK
}
