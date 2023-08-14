package cli

import (
	"io"
	"os"

	helpers "github.com/stutkhd-0709/enable_bootcamp/helpers"
)

func (f *InputFile) SplitBySize(sizePerFile int) error {
	buf := make([]byte, sizePerFile)

	fileCount := 0
	for {
		// buf: 読み込んだデータ
		// readByte: 読み込んだbyte数
		// readで読み込んだバイト数などの情報を持っているので毎回次のデータになる
		readByte, err := f.File.Read(buf)

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if readByte == 0 {
			break
		}

		outputFilename := helpers.GenerateFilename(f.NameWithoutExt, fileCount, f.Ext)
		err = os.WriteFile(outputFilename, buf[:readByte], 0644)
		if err != nil {
			return err
		}

		fileCount++
	}
	return nil
}
