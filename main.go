package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
)

var lineOpt int
var chunkOpt int
var sizeOpt int

func init() {
	// ポインタを指定して設定を予約
	flag.IntVar(&lineOpt, "l", 0, "分割ファイルの行数")
	flag.IntVar(&chunkOpt, "n", 0, "分割したいファイル数")
	flag.IntVar(&sizeOpt, "b", 0, "分割したいファイルサイズ")
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "[ERROR] length must be greater than 0, length = %d \n", flag.NArg())
		os.Exit(1)
	}

	if flag.NFlag() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	var filepath = flag.Args()[0]
	var splitErr error = nil
	if lineOpt > 0 {
		splitErr = splitFileByLine(filepath, lineOpt)
	} else if chunkOpt > 0 {
		splitErr = splitFileByChunk(filepath, chunkOpt)
	} else if sizeOpt > 0 {
		splitErr = splitFileBySize(filepath, sizeOpt)
	}

	if splitErr != nil {
		fmt.Fprintln(os.Stderr, splitErr)
		os.Exit(1)
	}

}

func splitFileByLine(filepath string, linesPerFile int) error {
	sf, err := os.Open(filepath)

	if err != nil {
		return err
	}

	// 関数終了時に閉じる
	defer sf.Close()

	scanner := bufio.NewScanner(sf)

	inputFilename := path.Base(filepath)
	inputFileExt := path.Ext(filepath)
	inputFileWithoutExt := inputFilename[:len(inputFilename)-len(inputFileExt)]

	var lineResult []byte
	lineCount := 0
	fileCount := 0

	for scanner.Scan() {
		line := append(scanner.Bytes(), byte('\n'))
		lineResult = append(lineResult, line...)
		lineCount++

		if lineCount%linesPerFile == 0 {
			outputFilename := generateFilename(inputFileWithoutExt, fileCount, inputFileExt)
			err := os.WriteFile(outputFilename, lineResult, 0644)
			if err != nil {
				return err
			}
			lineResult = []byte{}
			fileCount++
		}
	}

	// まとめてエラー処理をする
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "読み込みに失敗しました:", err)
	}
	return nil
}

func splitFileByChunk(filepath string, chunk int) error {
	return nil
}

func splitFileBySize(filepath string, size int) error {
	return nil
}

func generateFilename(prefix string, count int, extension string) string {
	characters := "abcdefghijklmnopqrstuvwxyz"
	suffix := ""
	for i := 0; i < 2; i++ {
		suffix = string(characters[count%26]) + suffix
		count = count / 26
	}
	return prefix + suffix + extension
}
