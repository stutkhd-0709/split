package cli

func (f *InputFile) SplitByChunk(fileChunk int, dist string) error {
	fileinfo, err := f.File.Stat()
	if err != nil {
		return err
	}

	fileSize := fileinfo.Size()

	chunkFileSize := int(fileSize) / fileChunk

	err = f.SplitBySize(chunkFileSize, dist)

	if err != nil {
		return err
	}

	return nil
}
