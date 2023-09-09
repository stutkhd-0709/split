package cli

import (
	"bufio"
	"io"
	"os"
	"sync"

	filehelpers "github.com/stutkhd-0709/split/filehelpers"
)

// filesizeはテストしやすいので外部から渡すようにする
type lineSplitter struct {
	reader       io.Reader
	divisionUnit int64
	dist         string
}

func NewLineSplitter(reader io.Reader, divisionUnit int64, dist string) Splitter {
	return &lineSplitter{reader, divisionUnit, dist}
}

func (s *lineSplitter) Split() (int64, error) {
	linesPerFile := s.divisionUnit
	scanner := bufio.NewScanner(s.reader)
	buf := make([]byte, 0, 64*1024)
	// 内部バッファをbufに置き換え、１行あたりの行数を拡大させる
	// これをループの内側に持ってくることは可能なのか？
	scanner.Buffer(buf, 1024*1024)

	var lineResult []byte
	var lineCount int64 = 0
	var fileCount int64 = 0

	var wg sync.WaitGroup
	errors := make(chan error)

	for scanner.Scan() {
		line := append(scanner.Bytes(), byte('\n'))
		lineResult = append(lineResult, line...)
		lineCount++
		if lineCount%linesPerFile == 0 {
			wg.Add(1)
			go func(_fileCount int64, _lineResult []byte, _dist string) {
				defer wg.Done()
				outputFilename, err := filehelpers.GenerateFilename(_dist, _fileCount)
				if err != nil {
					errors <- err
				}
				err = os.WriteFile(outputFilename, _lineResult, 0644)
				if err != nil {
					errors <- err
				}
			}(fileCount, lineResult, s.dist)
			fileCount++
			lineResult = []byte{}
		}
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
