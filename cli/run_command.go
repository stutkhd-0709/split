package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/stutkhd-0709/split/filehelpers"
)

type Splitter interface {
	Split() (int64, error)
}

type CLI struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

const (
	ExitOK int = 0
	ExitNG int = 1
)

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

func (cli *CLI) RunCommand(args []string) int {
	if err := ValidateArgs(); err != nil {
		return ExitNG
	}

	filepath := flag.Args()[0]

	f, fileSize, err := FileReader(filepath)
	if err != nil {
		return ExitNG
	}

	var dist string
	if len(flag.Args()) > 1 {
		dist = flag.Args()[1]
	} else {
		dist = ""
	}

	// flagを受け取って、optionを１つに指定する関数あるといいかも
	// func Option(f *flag) (Name, error)

	var splitter Splitter
	switch {
	case lineOpt != 0:
		splitter = NewLineSplitter(f, lineOpt, dist)
	case chunkOpt != 0:
		splitter = NewChunkSplitter(f, fileSize, chunkOpt, dist)
	case sizeOpt != "0":
		intSizeOpt, err := filehelpers.ConvertFileSizeToInt(sizeOpt)
		if err != nil {
			return ExitNG
		}
		splitter = NewSizeSplitter(f, fileSize, intSizeOpt, dist)
	default:
		return ExitNG
	}

	_, err = splitter.Split()

	if err != nil {
		return ExitNG
	}

	return ExitOK
}

func ValidateArgs() error {
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

	return nil

}

// getという名前は使用しない方がいいらしい
// https://google.github.io/styleguide/go/decisions#getters
func FileReader(filepath string) (io.Reader, int64, error) {
	_, err := os.Stat(filepath)
	if err != nil {
		return nil, 0, fmt.Errorf("ファイルが存在しません")
	}

	file, err := os.Open(filepath)

	if err != nil {
		return nil, 0, err
	}

	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}

	fileSize := fileinfo.Size()

	return file, fileSize, nil
}
