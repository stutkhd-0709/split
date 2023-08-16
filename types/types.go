package types

import (
	"io"
)

type InputFile struct {
	Reader   io.Reader
	FileSize int64
	FileName string
	Opt      *InputOpt
}

type InputOpt struct {
	LineOpt  int
	ChunkOpt int
	SizeOpt  string
}
