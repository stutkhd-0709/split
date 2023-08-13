package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
	"M": 1024 * 1024,
	"G": 1024 * 1024 * 1024,
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

	var wg sync.WaitGroup
	errors := make(chan error)

	for scanner.Scan() {
		line := append(scanner.Bytes(), byte('\n'))
		lineResult = append(lineResult, line...)
		lineCount++
		if lineCount%linesPerFile == 0 {
			wg.Add(1)
			go func(_fileCount int, _lineResult []byte) {
				defer wg.Done()
				outputFilename := generateFilename(f.NameWithoutExt, _fileCount, f.Ext)
				err := os.WriteFile(outputFilename, _lineResult, 0644)
				if err != nil {
					errors <- err
				}
			}(fileCount, lineResult)
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
			// Handle error
			fmt.Println("Error:", err)
		}
	}

	return nil
}

func (f *InputFile) splitFileBySize(sizePerFile int) error {
	buf := make([]byte, sizePerFile)

	fileCount := 0
	for {
		// buf: 読み込んだデータ
		// readByte: 読み込んだbyte数
		// readで読み込んだバイト数などの情報を持っているので毎回次のデータになる
		readByte, err := f.file.Read(buf)

		if err != nil {
			return err
		}
		if readByte == 0 {
			break
		}
		fileCount++

		outputFilename := generateFilename(f.NameWithoutExt, fileCount, f.Ext)
		err = os.WriteFile(outputFilename, buf[:readByte], 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *InputFile) splitFileByChunk(fileChunk int) error {
	fileinfo, err := f.file.Stat()
	if err != nil {
		return err
	}

	fileSize := fileinfo.Size()

	chunkFileSize := int(fileSize) / fileChunk

	err = f.splitFileBySize(chunkFileSize)

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

		sizeUnit := strings.ToUpper(matches[3])

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
	// rune型として扱う -> 文字コードのこと
	// 元の文字に戻すにはstring関数を使う
	firstChar := 'a' + (count % 26)
	secondChar := 'a' + (count / 26 % 26)

	fmt.Println(firstChar)
	fmt.Println(secondChar)
	// %cでUnicodeを表す
	return fmt.Sprintf("%s%c%c%s", prefix, firstChar, secondChar, extension)
}
