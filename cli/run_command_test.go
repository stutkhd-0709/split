package cli

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"testing"
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
		"noExistFile":              {"no_exits.txt", "", hasErr},
		"noOptionWithExistFile":    {"test.txt", "", hasErr},
		"overTwoArgument":          {"test.txt a b", "", hasErr},
		"LineOptionWithExistFile":  {"-l 2 test.txt ./test/", "", noErr},
		"sizeOptionWithExistFile":  {"-n 2 test.txt ./test/", "", noErr},
		"chunkOptionWithExistFile": {"-b 2 test.txt ./test/", "", noErr},
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

func TestSplit(t *testing.T) {
	t.Helper()
	t.Parallel()

	const (
		noErr  = false
		hasErr = true
	)

	// 空ファイルなので分割されない
	cases := map[string]struct {
		lengthPerline  int
		line           int
		lineOpt        int
		totalSplitFile int
		wantErr        bool
	}{
		"Split 2 files": {10, 2, 5, 2, noErr},
	}

	// テストを並列に実行する
	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data := make([]byte, tt.line*tt.lengthPerline)
			for i := range data {
				if i%tt.lengthPerline == tt.lengthPerline-1 {
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
					LineOpt:  tt.lineOpt,
					ChunkOpt: 0,
					SizeOpt:  "",
				},
			}

			fileCount, err := file.SplitByLine(tt.lineOpt, "./test/")

			switch {
			case tt.wantErr && err == nil:
				t.Fatal("expected error did not occur")
			case !tt.wantErr && err != nil:
				t.Fatal("unexpected error:", err)
			case fileCount != tt.totalSplitFile:
				t.Fatalf("Assertion Error: Get split file count is not expected, Get: %v, Expect: %v", fileCount, tt.totalSplitFile)
			}
		})
	}
}

// テストの前処理
func TestMain(m *testing.M) {
	file, err := os.Create("test.txt")
	if err != nil {
		log.Fatal(err)
	}

	status := m.Run()

	file.Close()

	if err := os.Remove("test.txt"); err != nil {
		log.Fatalf("Failed to remove test.txt: %v", err)
	}

	os.Exit(status)
}
