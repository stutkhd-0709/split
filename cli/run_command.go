package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/stutkhd-0709/split/model"
)

type CLI struct {
	// これいい感じにしたい
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

func (cli *CLI) RunCommand(args []string) error {
	var (
		lineOpt  int64
		chunkOpt int64
		sizeOpt  string
	)
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flagSet.Int64Var(&lineOpt, "l", 0, "分割ファイルの行数")
	flagSet.Int64Var(&chunkOpt, "n", 0, "分割したいファイル数")
	flagSet.StringVar(&sizeOpt, "b", "", "分割したいファイルサイズ")

	// テスト用にos.Argsを使用しないようにする
	err := flagSet.Parse(args[1:])

	if err != nil {
		return err
	}

	// TODO: validaterに切り分けたい
	switch {
	case flagSet.NArg() < 0:
		return fmt.Errorf("ファイルを指定してください")
	case flagSet.NArg() > 2:
		return fmt.Errorf("引数が多いです")
	case flagSet.NFlag() == 0:
		return fmt.Errorf("オプションを指定してください")
	case lineOpt == 0 && chunkOpt == 0 && sizeOpt == "":
		return fmt.Errorf("l, n, bのうちどれかオプションを指定してください")
	}

	filepath := flagSet.Args()[0]

	_, err = os.Stat(filepath)
	if err != nil {
		return fmt.Errorf("ファイルが存在しません")
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

	var splitter model.Splitter
	switch {
	case lineOpt != 0:
		splitter = NewLineSplitter(sf, fileSize, lineOpt)
	case chunkOpt != 0:
		splitter = NewChunkSplitter()
	case sizeOpt != "0":
		splitter = NewSizeSplitter()
	default:
		return fmt.Errorf("サイズを0以上に指定してください")
	}

	var dist string
	if len(flagSet.Args()) > 1 {
		dist = flagSet.Args()[1]
	} else {
		dist = ""
	}

	// こいつをモックにすることで、初期のファイルが不要になるかも
	_, err = splitter.Split(dist)

	if err != nil {
		return err
	}

	return nil
}
