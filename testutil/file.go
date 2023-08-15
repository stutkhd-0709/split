package testutil

import (
	"os"
	"syscall"
	"testing"
)

func CreateTestFile(t *testing.T, content []byte) string {
	file, err := os.CreateTemp("", "test")
	if err != nil {
		t.Error(err)
	}
	// ファイルを削除する
	t.Cleanup(func() { syscall.Unlink(file.Name()) })

	if err = os.WriteFile(file.Name(), content, 0644); err != nil {
		t.Error(err)
	}

	return file.Name()
}
