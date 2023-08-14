package cli

import (
	"flag"
	"fmt"
	"os"

	types "github.com/stutkhd-0709/enable_bootcamp/types"
)

const (
	ExitOK int = 0
	ExitNG int = 1
)

var lineOpt int
var chunkOpt int
var sizeOpt string

func init() {
	// ポインタを指定して設定を予約
	flag.IntVar(&lineOpt, "l", 0, "分割ファイルの行数")
	flag.IntVar(&chunkOpt, "n", 0, "分割したいファイル数")
	flag.StringVar(&sizeOpt, "b", "0", "分割したいファイルサイズ")
}

func Main() int {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "[ERROR] ファイルを指定してください")
		return ExitNG
	}

	if flag.NFlag() != 1 {
		fmt.Fprintln(os.Stderr, "[ERROR] オプションを指定してください")
		return ExitNG
	}

	filepath := flag.Args()[0]

	Opts := &types.InputOpt{
		LineOpt:  lineOpt,
		ChunkOpt: chunkOpt,
		SizeOpt:  sizeOpt,
	}

	cli := &CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	if err := cli.RunCommand(filepath, Opts); err != nil {
		fmt.Fprintln(cli.Stderr, "Error:", err)
		return ExitNG
	}

	return ExitOK
}
