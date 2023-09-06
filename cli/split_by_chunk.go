package cli

import (
	"fmt"
	"io"
	"os"
	"sync"

	filehelpers "github.com/stutkhd-0709/split/filehelpers"
	"github.com/stutkhd-0709/split/model"
)

type chunkSplitter struct {
	Reader       io.Reader
	FileSize     int64
	divisionUnit int64
}

func NewChunkSplitter(reader io.Reader, filesize int64, divisionUnit int64) model.Splitter {
	return &chunkSplitter{
		Reader:       reader,
		FileSize:     filesize,
		divisionUnit: divisionUnit,
	}
}

func (s *chunkSplitter) Split(dist string) (int64, error) {
	fileChunk := s.divisionUnit
	chunkFileSize := s.FileSize / fileChunk

	if chunkFileSize < 1 {
		return 0, fmt.Errorf("%s", "can't split into more than 0 files")
	}

	buf := make([]byte, chunkFileSize)

	var wg sync.WaitGroup
	errors := make(chan error)

	var fileCount int64 = 0
	var readByteSize int64 = 0
	for {
		// buf: 読み込んだデータ
		// readByte: 読み込んだbyte数
		// readで読み込んだバイト数などの情報を持っているので毎回次のデータになる
		if fileChunk-fileCount == 1 {
			buf = make([]byte, s.FileSize-readByteSize)
		}
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

		readByteSize += int64(readByte)

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
