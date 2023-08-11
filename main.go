package main

import (
	"flag"
	"fmt"
	"os"
)

var lineOpt int
var chunkOpt int
var sizeOpt int

func init() {
	// ポインタを指定して設定を予約
	flag.IntVar(&lineOpt, "l", 0, "分割ファイルの行数")
	flag.IntVar(&chunkOpt, "n", 0, "分割したいファイル数")
	flag.IntVar(&sizeOpt, "b", 0, "分割したいファイルサイズ")
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

	var filepath = flag.Args()[0]
	var err error = nil
	if lineOpt > 0 {
		err = SplitFileByLine(filepath, lineOpt)
	} else if chunkOpt > 0 {
		err = SplitFileByChunk(filepath, chunkOpt)
	} else if sizeOpt > 0 {
		err = SplitFileBySize(filepath, sizeOpt)
	}

	if err != nil {
		os.Exit(1)
	}

}

func SplitFileByLine(filepath string, line int) error {
	return nil
}

func SplitFileByChunk(filepath string, chunk int) error {
	return nil
}

func SplitFileBySize(filepath string, size int) error {
	return nil
}
