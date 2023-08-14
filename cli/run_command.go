package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	types "github.com/stutkhd-0709/enable_bootcamp/types"

	helpers "github.com/stutkhd-0709/enable_bootcamp/helpers"
)

type InputFile types.InputFile

type CLI struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

var lineOpt int
var chunkOpt int
var sizeOpt string

func init() {
	// ポインタを指定して設定を予約
	flag.IntVar(&lineOpt, "l", 0, "分割ファイルの行数")
	flag.IntVar(&chunkOpt, "n", 0, "分割したいファイル数")
	flag.StringVar(&sizeOpt, "b", "0", "分割したいファイルサイズ")
}

func (cli *CLI) RunCommand(args []string) error {
	flag.Parse()

	if flag.NArg() != 1 {
		return fmt.Errorf("[ERROR] ファイルを指定してください")
	}

	if flag.NFlag() != 1 {
		return fmt.Errorf("[ERROR] オプションを指定してください")
	}

	filepath := flag.Args()[0]

	Opts := &types.InputOpt{
		LineOpt:  lineOpt,
		ChunkOpt: chunkOpt,
		SizeOpt:  sizeOpt,
	}

	inputFilename := path.Base(filepath)
	inputFileExt := path.Ext(filepath)
	inputFileWithoutExt := inputFilename[:len(inputFilename)-len(inputFileExt)]

	sf, err := os.Open(filepath)

	if err != nil {
		return err
	}

	defer sf.Close()

	inputFile := &InputFile{
		File:           sf,
		NameWithoutExt: inputFileWithoutExt,
		Ext:            inputFileExt,
		Opt:            Opts,
	}

	if lineOpt != 0 {
		err = inputFile.SplitByLine(lineOpt)
	} else if chunkOpt != 0 {
		err = inputFile.SplitByChunk(chunkOpt)
	} else if sizeOpt != "0" {
		intFileSize, convertErr := helpers.ConvertFileSizeToInt(sizeOpt)
		if convertErr != nil {
			return err
		}
		err = inputFile.SplitBySize(intFileSize)
	}

	if err != nil {
		return err
	}

	return nil
}
