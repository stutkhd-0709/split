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

func (cli *CLI) RunCommand(args []string) int {
	// 内側に入れないと、並列テストの際に前後のテストに影響されてしまう
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
	err := flagSet.Parse(args)

	if err != nil {
		fmt.Println(err)
		return ExitNG
	}
	if err := ValidateArgs(flagSet); err != nil {
		fmt.Println(err)
		return ExitNG
	}

	filepath := flagSet.Args()[0]

	_, err = os.Stat(filepath)
	if err != nil {
		fmt.Println(err)
		return ExitNG
	}

	f, err := os.Open(filepath)

	if err != nil {
		fmt.Println(err)
		return ExitNG
	}

	fileinfo, err := f.Stat()
	if err != nil {
		fmt.Println(err)
		return ExitNG
	}

	fileSize := fileinfo.Size()

	if err != nil {
		fmt.Println(err)
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
	case sizeOpt != "":
		intSizeOpt, err := filehelpers.ConvertFileSizeToInt(sizeOpt)
		if err != nil {
			fmt.Println(err)
			return ExitNG
		}
		splitter = NewSizeSplitter(f, fileSize, intSizeOpt, dist)
	default:
		fmt.Println("オプションを指定してください")
		return ExitNG
	}

	_, err = splitter.Split()

	if err != nil {
		fmt.Println(err)
		return ExitNG
	}

	return ExitOK
}

func ValidateArgs(fs *flag.FlagSet) error {
	switch {
	case fs.NArg() < 0:
		return fmt.Errorf("ファイルを指定してください")
	case fs.NArg() > 2:
		return fmt.Errorf("引数が多いです")
	case fs.NFlag() == 0:
		return fmt.Errorf("オプションを指定してください")
	}

	return nil

}
