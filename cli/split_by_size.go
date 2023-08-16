package cli

import (
	"io"
	"os"
	"sync"

	helpers "github.com/stutkhd-0709/enable_bootcamp/helpers"
)

func (f *InputFile) SplitBySize(sizePerFile int, dist string) error {
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
			return err
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
			return err
		}
	}

	return nil
}
