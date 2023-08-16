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

func (cli *CLI) RunCommand(args []string) error {
	var (
		lineOpt  int
		chunkOpt int
		sizeOpt  string
	)
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flagSet.IntVar(&lineOpt, "l", 0, "分割ファイルの行数")
	flagSet.IntVar(&chunkOpt, "n", 0, "分割したいファイル数")
	flagSet.StringVar(&sizeOpt, "b", "", "分割したいファイルサイズ")

	// テスト用にos.Argsを使用しないようにする
	err := flagSet.Parse(args[1:])

	if err != nil {
		return err
	}

	if flagSet.NArg() < 0 {
		return fmt.Errorf("ファイルを指定してください")
	}

	if flagSet.NArg() > 2 {
		return fmt.Errorf("引数が多いです")
	}

	if flagSet.NFlag() == 0 {
		return fmt.Errorf("オプションを指定してください")
	}

	if lineOpt == 0 && chunkOpt == 0 && sizeOpt == "" {
		return fmt.Errorf("l, n, bのうちどれかオプションを指定してください")
	}

	filepath := flagSet.Args()[0]

	var dist string
	if len(flagSet.Args()) > 1 {
		dist = flagSet.Args()[1]
	} else {
		dist = ""
	}

	_, err = os.Stat(filepath)
	if err != nil {
		return fmt.Errorf("ファイルが存在しません")
	}

	Opts := &types.InputOpt{
		LineOpt:  lineOpt,
		ChunkOpt: chunkOpt,
		SizeOpt:  sizeOpt,
	}

	sf, err := os.Open(filepath)

	if err != nil {
		return err
	}

	defer sf.Close()

	fileinfo, err := sf.Stat()
	if err != nil {
		return err
	}

	fileSize := fileinfo.Size()

	inputFile := &InputFile{
		Reader:   sf,
		FileSize: fileSize,
		FileName: path.Base(filepath),
		Opt:      Opts,
	}

	if lineOpt != 0 {
		err = inputFile.SplitByLine(lineOpt, dist)
	} else if chunkOpt != 0 {
		err = inputFile.SplitByChunk(chunkOpt, dist)
	} else if sizeOpt != "0" {
		intFileSize, convertErr := helpers.ConvertFileSizeToInt(sizeOpt)
		if convertErr != nil {
			return err
		}
		err = inputFile.SplitBySize(intFileSize, dist)
	}

	if err != nil {
		return err
	}

	return nil
}
