package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	helpers "github.com/stutkhd-0709/enable_bootcamp/helpers"
)

type InputFile struct {
	Reader   io.Reader
	FileSize int64
	FileName string
	Opt      *InputOpt
}

type InputOpt struct {
	LineOpt  int
	ChunkOpt int
	SizeOpt  string
}

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

	Opts := &InputOpt{
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
		_, err = inputFile.SplitByLine(lineOpt, dist)
	} else if chunkOpt != 0 {
		_, err = inputFile.SplitByChunk(chunkOpt, dist)
	} else if sizeOpt != "0" {
		intFileSize, convertErr := helpers.ConvertFileSizeToInt(sizeOpt)
		if convertErr != nil {
			return err
		}
		_, err = inputFile.SplitBySize(intFileSize, dist)
	}

	if err != nil {
		return err
	}

	return nil
}

func (f *InputFile) SplitByChunk(fileChunk int, dist string) (int, error) {
	chunkFileSize := int(f.FileSize) / fileChunk

	fileCount, err := f.SplitBySize(chunkFileSize, dist)

	if err != nil {
		return 0, err
	}

	return fileCount, nil
}

func (f *InputFile) SplitBySize(sizePerFile int, dist string) (int, error) {
	buf := make([]byte, sizePerFile)

	var wg sync.WaitGroup
	errors := make(chan error)

	fileCount := 0
	for {
		// buf: 読み込んだデータ
		// readByte: 読み込んだbyte数
		// readで読み込んだバイト数などの情報を持っているので毎回次のデータになる
		readByte, err := f.Reader.Read(buf)

		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
		if readByte == 0 {
			break
		}

		wg.Add(1)
		writeBuf := make([]byte, readByte)
		copy(writeBuf, buf[:readByte])
		go func(_fileCount int, _writBuf []byte, _dist string) {
			defer wg.Done()
			outputFilename, err := helpers.GenerateFilename(_dist, _fileCount)
			if err != nil {
				errors <- err
			}
			err = os.WriteFile(outputFilename, _writBuf, 0644)
			if err != nil {
				errors <- err
			}
		}(fileCount, writeBuf, dist)

		fileCount++
	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	for err := range errors {
		if err != nil {
			return 0, err
		}
	}

	return fileCount, nil
}

func (f *InputFile) SplitByLine(linesPerFile int, dist string) (int, error) {
	scanner := bufio.NewScanner(f.Reader)

	var lineResult []byte
	lineCount := 0
	fileCount := 0

	var wg sync.WaitGroup
	errors := make(chan error)

	for scanner.Scan() {
		line := append(scanner.Bytes(), byte('\n'))
		lineResult = append(lineResult, line...)
		lineCount++
		if lineCount%linesPerFile == 0 {
			wg.Add(1)
			go func(_fileCount int, _lineResult []byte, _dist string) {
				defer wg.Done()
				outputFilename, err := helpers.GenerateFilename(_dist, _fileCount)
				if err != nil {
					errors <- err
				}
				err = os.WriteFile(outputFilename, _lineResult, 0644)
				if err != nil {
					errors <- err
				}
			}(fileCount, lineResult, dist)
			fileCount++
			lineResult = []byte{}
		}
	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	// エラーがあるか確認
	// 上のgoroutineが実行されるまで行う
	// 子goroutineが処理を終了するたびに実行される
	// for rangeをchannelで行う場合、goではそのチャンネルがcloseされるまで実行される
	for err := range errors {
		if err != nil {
			return 0, err
		}
	}

	return fileCount, nil
}
