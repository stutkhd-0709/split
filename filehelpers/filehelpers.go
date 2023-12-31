package filehelpers

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var fileSizeUnitToBytes = map[string]int{
	"K": 1024,
	"M": 1024 * 1024,
	"G": 1024 * 1024 * 1024,
}

// 別のファイルに分けても良さそう -> byteの時しか使用しないので
func ConvertFileSizeToInt(strFileSize string) (int64, error) {
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

	return int64(intFileByteSize), nil
}

func GenerateFilename(dist string, count int64) (string, error) {
	// rune型として扱う -> 文字コードのこと
	// 元の文字に戻すにはstring関数を使う
	firstChar := 'a' + (count / 26 % 26)
	secondChar := 'a' + (count % 26)

	var prefix string
	if dist == "" {
		prefix = "x"
	} else {
		isDir, err := ensureDirectoryExists(dist)
		if err != nil {
			return "", err
		}

		if isDir {
			prefix = dist + "x"
		} else {
			prefix = dist
		}
	}

	// %cでUnicodeを表す
	return fmt.Sprintf("%s%c%c", prefix, firstChar, secondChar), nil
}

func ensureDirectoryExists(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		// pathが存在しない場合
		if os.IsNotExist(err) {
			// ディレクトリかどうか
			if strings.HasSuffix(path, "/") {
				// 新規作成
				err := os.MkdirAll(path, 0755)
				if err != nil {
					return true, err
				}
				return true, nil
			} else {
				return false, nil
			}
		} else {
			return false, err
		}
	} else if !fileInfo.IsDir() {
		// pathが存在して、そのpathがディレクトリでない場合
		return false, nil
	}
	// すでにpathが存在している場合
	return true, nil
}
