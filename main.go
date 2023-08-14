package main

import (
	"flag"
	"fmt"
	"os"

	cli "github.com/stutkhd-0709/enable_bootcamp/cli"
	types "github.com/stutkhd-0709/enable_bootcamp/types"
)

var lineOpt int
var chunkOpt int
var sizeOpt string

var fileSizeUnitToBytes = map[string]int{
	"K": 1024,
	"M": 1024 * 1024,
	"G": 1024 * 1024 * 1024,
}

func init() {
	// ポインタを指定して設定を予約
	flag.IntVar(&lineOpt, "l", 0, "分割ファイルの行数")
	flag.IntVar(&chunkOpt, "n", 0, "分割したいファイル数")
	flag.StringVar(&sizeOpt, "b", "0", "分割したいファイルサイズ")
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "[ERROR] length must be greater than 0, length = %d \n", flag.NArg())
		os.Exit(1)
	}

	if flag.NFlag() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	filepath := flag.Args()[0]

	Opts := &types.InputOpt{
		LineOpt:  lineOpt,
		ChunkOpt: chunkOpt,
		SizeOpt:  sizeOpt,
	}

	err := cli.RunSplit(filepath, Opts)

	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v \n", err)
		os.Exit(1)
	}
}
