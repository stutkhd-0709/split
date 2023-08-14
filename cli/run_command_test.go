package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	t.Helper()
	t.Parallel()

	const (
		noErr  = false
		hasErr = true
	)

	cases := map[string]struct {
		args    string
		in      string
		wantErr bool
	}{
		"noArgument": {"", "", hasErr},
	}

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

			args := strings.Split(tt.args, " ")
			os.Args = append([]string{"cmd"}, args...)

			err := cli.RunCommand(os.Args)

			switch {
			case tt.wantErr && err == nil:
				t.Fatal("expected error did not occur")
			case !tt.wantErr && err != nil:
				t.Fatal("unexpected error:", err)
			}
		})
	}
}
