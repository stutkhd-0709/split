package types

import (
	"os"
)

type InputFile struct {
	File           *os.File
	NameWithoutExt string
	Ext            string
	Opt            *InputOpt
}

type InputOpt struct {
	LineOpt  int
	ChunkOpt int
	SizeOpt  string
}
