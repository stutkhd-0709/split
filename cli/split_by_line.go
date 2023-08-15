package cli

import (
	"bufio"
	"os"
	"sync"

	helpers "github.com/stutkhd-0709/enable_bootcamp/helpers"
)

func (f *InputFile) SplitByLine(linesPerFile int, dist string) error {
	scanner := bufio.NewScanner(f.File)

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
			return err
		}
	}

	return nil
}
