package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
)

var lineOpt int
var chunkOpt int
var sizeOpt string

type InputFile struct {
	file           *os.File
	NameWithoutExt string
	Ext            string
}

var fileSizeUnitToBytes = map[string]int{
	"K": 1024,
	"k": 1024,
	"M": 1024 * 1024,
	"m": 1024 * 1024,
	"G": 1024 * 1024 * 1024,
	"g": 1024 * 1024 * 1024,
}

func init() {
	// ポインタを指定して設定を予約
	flag.IntVar(&lineOpt, "l", 0, "分割ファイルの行数")
	flag.IntVar(&chunkOpt, "n", 0, "分割したいファイル数")
	flag.StringVar(&sizeOpt, "b", "0", "分割したいファイルサイズ")
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

	err := runSplit(filepath)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runSplit(filepath string) error {
	inputFilename := path.Base(filepath)
	inputFileExt := path.Ext(filepath)
	inputFileWithoutExt := inputFilename[:len(inputFilename)-len(inputFileExt)]

	sf, err := os.Open(filepath)

	if err != nil {
		return err
	}

	defer sf.Close()

	inputFile := &InputFile{
		file:           sf,
		NameWithoutExt: inputFileWithoutExt,
		Ext:            inputFileExt,
	}

	if lineOpt != 0 {
		err = inputFile.splitFileByLine(lineOpt)
	} else if chunkOpt != 0 {
		err = inputFile.splitFileByChunk(chunkOpt)
	} else if sizeOpt != "0" {
		intFileSize, convertErr := convertFileSizeToInt(sizeOpt)
		if convertErr != nil {
			return err
		}
		err = inputFile.splitFileBySize(intFileSize)
	}

	if err != nil {
		return err
	}

	return nil
}

func (f *InputFile) splitFileByLine(linesPerFile int) error {
	scanner := bufio.NewScanner(f.file)

	var lineResult []byte
	lineCount := 0
	fileCount := 0

	for scanner.Scan() {
		line := append(scanner.Bytes(), byte('\n'))
		lineResult = append(lineResult, line...)
		lineCount++

		if lineCount%linesPerFile == 0 {
			outputFilename := generateFilename(f.NameWithoutExt, fileCount, f.Ext)
			err := os.WriteFile(outputFilename, lineResult, 0644)
			if err != nil {
				return err
			}
			lineResult = []byte{}
			fileCount++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "読み込みに失敗しました:", err)
	}
	return nil
}

func (f *InputFile) splitFileBySize(size int) error {
	fmt.Println(size)
	return nil
}

func (f *InputFile) splitFileByChunk(fileChunk int) error {
	fileinfo, err := f.file.Stat()
	if err != nil {
		return err
	}

	fileSize := fileinfo.Size()

	chunkFileSize := fileSize / int64(fileChunk)

	err = f.splitFileBySize(int(chunkFileSize))

	if err != nil {
		return err
	}

	return nil
}

func convertFileSizeToInt(strFileSize string) (int, error) {
	numericPattern := `^\d+$`
	match, err := regexp.MatchString(numericPattern, strFileSize)

	if err != nil {
		return 0, err
	}

	var intFileByteSize int
	if match {
		intFileByteSize, _ = strconv.Atoi(strFileSize)
	} else {
		re := regexp.MustCompile(`^(\d+(\.\d+)?)([KkMmGg])$`)
		matches := re.FindStringSubmatch(strFileSize)

		if len(matches) == 0 {
			argumentErr := fmt.Errorf("error: サイズの指定方法が間違えています")
			return 0, argumentErr
		}

		inputFileSize, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, err
		}

		sizeUnit := matches[3]

		unitToByte := fileSizeUnitToBytes[sizeUnit]
		if unitToByte == 0 {
			argumentErr := fmt.Errorf("error: 対応してないサイズ単位です")
			return 0, argumentErr
		}

		intFileByteSize = inputFileSize * unitToByte
	}

	return intFileByteSize, nil
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
