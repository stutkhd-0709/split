package cli

import (
	"bytes"
	"io"
	"testing"

	"github.com/stutkhd-0709/enable_bootcamp/types"
)

func TestSplitByLine(t *testing.T) {
	t.Helper()
	t.Parallel()

	const (
		noErr  = false
		hasErr = true
	)

	// 空ファイルなので分割されない
	cases := map[string]struct {
		lengthPerline int
		line          int
		lineOpt       int
		wantErr       bool
	}{}

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
				Opt: &types.InputOpt{
					LineOpt:  tt.lineOpt,
					ChunkOpt: 0,
					SizeOpt:  "",
				},
			}

			err := file.SplitByLine(tt.lineOpt, "./test/")

			switch {
			case tt.wantErr && err == nil:
				t.Fatal("expected error did not occur")
			case !tt.wantErr && err != nil:
				t.Fatal("unexpected error:", err)
			}
		})
	}
}
