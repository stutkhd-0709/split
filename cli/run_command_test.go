package cli

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestRunCommand(t *testing.T) {
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

			// グローバル変数だよねん
			// https://stackoverflow.com/questions/33723300/how-to-test-the-passing-of-arguments-in-golang
			// deferで元の状態に戻すとかやってるけど、それ以外にいい方法あるのかどうか
			os.Args = []string{"cmd", "-l", "2", "./test/test.txt", "./test/"}

			fmt.Println(os.Args)

			// テストするときにos.Argsに値を入れるのは良くない？
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
