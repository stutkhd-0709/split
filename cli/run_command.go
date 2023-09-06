package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/stutkhd-0709/split/filehelpers"
	"github.com/stutkhd-0709/split/model"
)

type CLI struct {
	// これいい感じにしたい
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

var (
	lineOpt  int64
	chunkOpt int64
	sizeOpt  string
)

func init() {
	flag.Int64Var(&lineOpt, "l", 0, "分割ファイルの行数")
	flag.Int64Var(&chunkOpt, "n", 0, "分割したいファイル数")
	flag.StringVar(&sizeOpt, "b", "", "分割したいファイルサイズ")
}

func (cli *CLI) RunCommand(args []string) error {
	filepath := flag.Args()[0]

	if err := validateArgs(filepath); err != nil {
		return fmt.Errorf(err.Error())
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
		splitter = NewChunkSplitter(sf, fileSize, chunkOpt)
	case sizeOpt != "0":
		intSizeOpt, err := filehelpers.ConvertFileSizeToInt(sizeOpt)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		splitter = NewSizeSplitter(sf, fileSize, intSizeOpt)
	default:
		return fmt.Errorf("サイズを0以上に指定してください")
	}

	var dist string
	if len(flag.Args()) > 1 {
		dist = flag.Args()[1]
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

func validateArgs(filepath string) error {
	flag.Parse()
	switch {
	case flag.NArg() < 0:
		return fmt.Errorf("ファイルを指定してください")
	case flag.NArg() > 2:
		return fmt.Errorf("引数が多いです")
	case flag.NFlag() == 0:
		return fmt.Errorf("オプションを指定してください")
	case lineOpt == 0 && chunkOpt == 0 && sizeOpt == "":
		return fmt.Errorf("l, n, bのうちどれかオプションを指定してください")
	}

	_, err := os.Stat(filepath)
	if err != nil {
		return fmt.Errorf("ファイルが存在しません")
	}

	return nil
}
