package cli

func (f *InputFile) SplitByChunk(fileChunk int, dist string) error {
	chunkFileSize := int(f.FileSize) / fileChunk

	err := f.SplitBySize(chunkFileSize, dist)

	if err != nil {
		return err
	}

	return nil
}
