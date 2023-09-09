package cli

import (
	"fmt"
	"io"
	"os"
	"sync"

	filehelpers "github.com/stutkhd-0709/split/filehelpers"
)

type sizeSplitter struct {
	Reader       io.Reader
	FileSize     int64
	divisionUnit int64
	dist         string
}

func NewSizeSplitter(reader io.Reader, filesize int64, divisionUnit int64, dist string) Splitter {
	return &sizeSplitter{
		Reader:       reader,
		FileSize:     filesize,
		divisionUnit: divisionUnit,
		dist:         dist,
	}
}

func (s *sizeSplitter) Split() (int64, error) {
	sizePerFile := s.divisionUnit

	if sizePerFile > s.FileSize {
		return 0, fmt.Errorf("%s", "Specified size exceeds original file size")
	}

	buf := make([]byte, sizePerFile)

	var wg sync.WaitGroup
	errors := make(chan error)

	var fileCount int64 = 0
	for {
		// buf: 読み込んだデータ
		// readByte: 読み込んだbyte数
		// readで読み込んだバイト数などの情報を持っているので毎回次のデータになる
		readByte, err := s.Reader.Read(buf)

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
		go func(_fileCount int64, _writBuf []byte, _dist string) {
			defer wg.Done()
			outputFilename, err := filehelpers.GenerateFilename(_dist, _fileCount)
			if err != nil {
				errors <- err
			}
			err = os.WriteFile(outputFilename, _writBuf, 0644)
			if err != nil {
				errors <- err
			}
		}(fileCount, writeBuf, s.dist)

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
