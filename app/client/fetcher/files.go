package fetcher

import (
	"errors"
	"os"
)

func (a *FileFetcher) getFile(numBytes int, filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return a.newFile(numBytes, filename)
	}

	info, _ := file.Stat()
	if info.Size() == int64(numBytes) {
		return file, nil
	}

	return a.newFile(numBytes, filename)
}

func (a *FileFetcher) newFile(numBytes int, filename string) (*os.File, error) {
	file, err := os.Create(filename)

	handleError(err)

	bytes := make([]byte, numBytes)
	n, err := file.Write(bytes)

	handleError(err)

	if n != numBytes {
		file.Close()
		return nil, errors.New("bytes written and deisred size do not match")
	}

	return file, nil
}
