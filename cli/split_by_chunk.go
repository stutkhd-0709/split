package cli

func (f *InputFile) SplitByChunk(fileChunk int) error {
	fileinfo, err := f.File.Stat()
	if err != nil {
		return err
	}

	fileSize := fileinfo.Size()

	chunkFileSize := int(fileSize) / fileChunk

	err = f.SplitBySize(chunkFileSize)

	if err != nil {
		return err
	}

	return nil
}
