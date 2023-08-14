package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var fileSizeUnitToBytes = map[string]int{
	"K": 1024,
	"M": 1024 * 1024,
	"G": 1024 * 1024 * 1024,
}

func ConvertFileSizeToInt(strFileSize string) (int, error) {
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

func GenerateFilename(prefix string, count int, extension string) string {
	// rune型として扱う -> 文字コードのこと
	// 元の文字に戻すにはstring関数を使う
	firstChar := 'a' + (count % 26)
	secondChar := 'a' + (count / 26 % 26)

	// %cでUnicodeを表す
	return fmt.Sprintf("%s%c%c%s", prefix, firstChar, secondChar, extension)
}
