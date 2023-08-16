package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	helpers "github.com/stutkhd-0709/enable_bootcamp/helpers"
)

func TestRunCommand(t *testing.T) {
	t.Helper()
	t.Parallel()

	const (
		noErr  = false
		hasErr = true
	)

	// 空ファイルなので分割されない
	cases := map[string]struct {
		args    string
		in      string
		wantErr bool
	}{
		"noArgument":               {"", "", hasErr},
		"noExistFile":              {"./test/no_exits.txt", "", hasErr},
		"noOptionWithExistFile":    {"./test/test.txt", "", hasErr},
		"overTwoArgument":          {"./test/test.txt a b", "", hasErr},
		"LineOptionWithExistFile":  {"-l 2 ./test/test.txt ./test/", "", noErr},
		"chunkOptionWithExistFile": {"-n 2 ./test/test.txt ./test/", "", noErr},
		"sizeOptionWithExistFile":  {"-b 2 ./test/test.txt ./test/", "", noErr},
	}

	// テストを並列に実行する
	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var got bytes.Buffer
			cli := &CLI{
				Stdout: &got,
				Stderr: &got,
				Stdin:  strings.NewReader(tt.in),
			}

			var args []string
			if tt.args != "" {
				args = strings.Split(tt.args, " ")
			} else {
				args = []string{}
			}

			args = append([]string{"cmd"}, args...)

			err := cli.RunCommand(args)

			switch {
			case tt.wantErr && err == nil:
				t.Fatal("expected error did not occur")
			case !tt.wantErr && err != nil:
				t.Fatal("unexpected error:", err)
			}
		})
	}
}

func TestSplitByLine(t *testing.T) {
	t.Helper()
	t.Parallel()

	const (
		noErr  = false
		hasErr = true
	)

	type lineOpts struct {
		perLine   int64
		line      int64
		splitLine int64
	}

	cases := map[string]struct {
		wantErr     bool
		expectSplit int64
		lineOpts    lineOpts
	}{
		"Split 5 files": {
			noErr,
			5,
			lineOpts{
				perLine:   2,
				line:      10,
				splitLine: 2,
			},
		},
		"split perLine 10000 file to two": {
			noErr,
			2,
			lineOpts{
				perLine:   10000,
				line:      2,
				splitLine: 1,
			},
		},
		"split perLine 100000 file to two": {
			noErr,
			2,
			lineOpts{
				perLine:   100000,
				line:      2,
				splitLine: 1,
			},
		},
		"split perLine 1024 * 1024 file to two file": {
			noErr,
			2,
			lineOpts{
				perLine:   1024 * 1024,
				line:      2,
				splitLine: 1,
			},
		},
		"split 1000 lines file to 1000 file": {
			noErr,
			1000,
			lineOpts{
				perLine:   1,
				line:      1000,
				splitLine: 1,
			},
		},
	}

	// テストを並列に実行する
	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			lineOpts := tt.lineOpts
			data := make([]byte, lineOpts.line*lineOpts.perLine)
			for i := range data {
				if int64(i)%lineOpts.perLine == lineOpts.perLine-1 {
					data[i] = '\n'
				} else {
					data[i] = 'a'
				}
			}

			var r io.Reader = bytes.NewBuffer(data)

			file := &InputFile{
				Reader:   r,
				FileName: "test",
				Opt: &InputOpt{
					LineOpt:  lineOpts.splitLine,
					ChunkOpt: 0,
					SizeOpt:  "",
				},
			}

			dist := "./test/" + t.Name() + "/"
			fileCount, err := file.SplitByLine(dist)

			defer func() {
				err := os.RemoveAll(dist)
				if err != nil {
					fmt.Println("Failed to delete files:", err)
				}
			}()

			switch {
			case tt.wantErr && err == nil:
				t.Fatal("expected error did not occur")
			case !tt.wantErr && err != nil:
				t.Fatal("unexpected error:", err)
			case fileCount != tt.expectSplit:
				t.Fatalf("Assertion Error: Get split file count is not expected, Get: %v, Expect: %v", fileCount, tt.expectSplit)
			}
		})
	}
}

func TestSplitBySize(t *testing.T) {
	t.Helper()
	t.Parallel()

	const (
		noErr  = false
		hasErr = true
	)

	type sizeOpts struct {
		splitSizeWithUnit string
		fileSizeWithUnit  string
	}

	cases := map[string]struct {
		wantErr     bool
		expectSplit int64
		sizeOpts    sizeOpts
	}{
		"Specified Split Size Over Original": {
			hasErr,
			0,
			sizeOpts{
				splitSizeWithUnit: "10K",
				fileSizeWithUnit:  "1k",
			},
		},
		"Split a 100k file to 5 files": {
			noErr,
			10,
			sizeOpts{
				splitSizeWithUnit: "10K",
				fileSizeWithUnit:  "100k",
			},
		},
		"Split a 100M files to 10 files": {
			noErr,
			10,
			sizeOpts{
				splitSizeWithUnit: "10M",
				fileSizeWithUnit:  "100M",
			},
		},
		"Split a 1G files to 2 files": {
			noErr,
			2,
			sizeOpts{
				splitSizeWithUnit: "512M",
				fileSizeWithUnit:  "1G",
			},
		},
		"Split a 10K files to 4 files": {
			noErr,
			4,
			sizeOpts{
				splitSizeWithUnit: "3K",
				fileSizeWithUnit:  "10K",
			},
		},
		"Split a 100K files to 100 files": {
			noErr,
			100,
			sizeOpts{
				splitSizeWithUnit: "1K",
				fileSizeWithUnit:  "100K",
			},
		},
	}

	// テストを並列に実行する
	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			sizeOpts := tt.sizeOpts

			sizeOriginFile, err := helpers.ConvertFileSizeToInt(sizeOpts.fileSizeWithUnit)

			if err != nil {
				t.Fatal("convert origin filesize to int:", err)
			}

			data := make([]byte, sizeOriginFile)
			for i := range data {
				data[i] = 'a'
			}

			var r io.Reader = bytes.NewBuffer(data)

			file := &InputFile{
				Reader:   r,
				FileName: "test",
				FileSize: sizeOriginFile,
				Opt: &InputOpt{
					LineOpt:  0,
					ChunkOpt: 0,
					SizeOpt:  sizeOpts.splitSizeWithUnit,
				},
			}

			dist := "./test/" + t.Name() + "/"
			fileCount, err := file.SplitBySize(dist)

			defer func() {
				err := os.RemoveAll(dist)
				if err != nil {
					fmt.Println("Failed to delete files:", err)
				}
			}()

			switch {
			case tt.wantErr && err == nil:
				t.Fatal("expected error did not occur")
			case !tt.wantErr && err != nil:
				t.Fatal("unexpected error:", err)
			case fileCount != tt.expectSplit:
				t.Fatalf("Assertion Error: Get split file count is not expected, Get: %v, Expect: %v", fileCount, tt.expectSplit)
			}
		})
	}
}

func TestSplitByChunk(t *testing.T) {
	t.Helper()
	t.Parallel()

	const (
		noErr  = false
		hasErr = true
	)

	type chunkOpts struct {
		chunkSize        int64
		fileSizeWithUnit string
	}

	cases := map[string]struct {
		wantErr     bool
		expectSplit int64
		chunkOpts   *chunkOpts
	}{
		"Split No Size File": {
			hasErr,
			0,
			&chunkOpts{
				chunkSize:        3,
				fileSizeWithUnit: "0",
			},
		},
		"Split Odd Chunk": {
			noErr,
			3,
			&chunkOpts{
				chunkSize:        3,
				fileSizeWithUnit: "10K",
			},
		},
		"Split 10k to 2 files": {
			noErr,
			2,
			&chunkOpts{
				chunkSize:        2,
				fileSizeWithUnit: "10K",
			},
		},
		"Split 100M to 2 files": {
			noErr,
			2,
			&chunkOpts{
				chunkSize:        2,
				fileSizeWithUnit: "100M",
			},
		},
		"Split 1G to 2 files": {
			noErr,
			2,
			&chunkOpts{
				chunkSize:        2,
				fileSizeWithUnit: "1G",
			},
		},
	}

	// テストを並列に実行する
	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			chunkOpts := tt.chunkOpts

			sizeOriginFile, err := helpers.ConvertFileSizeToInt(chunkOpts.fileSizeWithUnit)

			if err != nil {
				t.Fatal("convert origin filesize to int:", err)
			}

			data := make([]byte, sizeOriginFile)
			for i := range data {
				data[i] = 'a'
			}

			var r io.Reader = bytes.NewBuffer(data)

			file := &InputFile{
				Reader:   r,
				FileName: "test",
				FileSize: sizeOriginFile,
				Opt: &InputOpt{
					LineOpt:  0,
					ChunkOpt: chunkOpts.chunkSize,
					SizeOpt:  "",
				},
			}

			dist := "./test/" + t.Name() + "/"
			fileCount, err := file.SplitByChunk(dist)

			defer func() {
				err := os.RemoveAll(dist)
				if err != nil {
					fmt.Println("Failed to delete files:", err)
				}
			}()

			switch {
			case tt.wantErr && err == nil:
				t.Fatal("expected error did not occur")
			case !tt.wantErr && err != nil:
				t.Fatal("unexpected error:", err)
			case fileCount != tt.expectSplit:
				t.Fatalf("Assertion Error: Get split file count is not expected, Get: %v, Expect: %v", fileCount, tt.expectSplit)
			}
		})
	}
}

// テストの前処理
func TestMain(m *testing.M) {
	dir := "./test"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.Create(dir + "/test.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := []string{"a\n", "b\n", "c\n"}
	for _, line := range lines {
		b := []byte(line)
		_, err := file.Write(b)
		if err != nil {
			log.Fatal(err)
		}
	}

	m.Run()

	defer file.Close()

	if err := os.Remove(dir + "/test.txt"); err != nil {
		log.Fatalf("Failed to remove test.txt: %v", err)
	}

	err = os.RemoveAll("./test/")
	if err != nil {
		fmt.Println("Failed to delete files:", err)
	}
}
