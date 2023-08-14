package cli

import (
	"os"
	"path"

	types "github.com/stutkhd-0709/enable_bootcamp/types"

	helpers "github.com/stutkhd-0709/enable_bootcamp/helpers"
)

type InputFile types.InputFile

func RunSplit(filepath string, inputOpt *types.InputOpt) error {
	inputFilename := path.Base(filepath)
	inputFileExt := path.Ext(filepath)
	inputFileWithoutExt := inputFilename[:len(inputFilename)-len(inputFileExt)]

	sf, err := os.Open(filepath)

	if err != nil {
		return err
	}

	defer sf.Close()

	inputFile := &InputFile{
		File:           sf,
		NameWithoutExt: inputFileWithoutExt,
		Ext:            inputFileExt,
		Opt:            inputOpt,
	}

	if inputOpt.LineOpt != 0 {
		err = inputFile.SplitByLine(inputOpt.LineOpt)
	} else if inputOpt.ChunkOpt != 0 {
		err = inputFile.SplitByChunk(inputOpt.ChunkOpt)
	} else if inputOpt.SizeOpt != "0" {
		intFileSize, convertErr := helpers.ConvertFileSizeToInt(inputOpt.SizeOpt)
		if convertErr != nil {
			return err
		}
		err = inputFile.SplitBySize(intFileSize)
	}

	if err != nil {
		return err
	}

	return nil
}
