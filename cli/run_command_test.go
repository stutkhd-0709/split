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
		args    []string
		in      string
		wantErr bool
	}{
		// "引数が何もない": {[]string{""}, "", hasErr},
		// "存在しないファイルを指定": {[]string{"./test/no_exits.txt"}, "", hasErr},
		// こいつが悪さしてそう
		"存在するファイルにオプション指定なし": {[]string{"./test/test.txt"}, "", hasErr},
		"引数を２つ以上指定":          {[]string{"./test/test.txt", "a", "b"}, "", hasErr},
		"存在するファイルにライン分割指定":   {[]string{"-l", "2", "./test/test.txt", "./test/"}, "", noErr},
		// "存在するファイルにチャンク分割指定":  {[]string{"-n", "2", "./test/test.txt", "./test/"}, "", noErr},
		// "存在するファイルにサイズ分割指定":   {[]string{"-b", "2", "./test/test.txt", "./test/"}, "", noErr},
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

			exitCode := cli.RunCommand(tt.args)

			switch {
			case tt.wantErr && exitCode == 0:
				t.Fatal("expected error did not occur")
			case !tt.wantErr && exitCode != 0:
				t.Fatal("unexpected exitCode:", exitCode)
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
