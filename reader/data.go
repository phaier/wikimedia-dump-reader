package reader

import (
	"bytes"
	"compress/bzip2"
	"io"
	"os"
)

func readStream(file *os.File, position *Index) (string, error) {
	_, err := file.Seek(position.Offset, io.SeekStart)
	if err != nil {
		return "", err
	}

	source := make([]byte, position.Length)
	_, err = io.ReadFull(file, source)
	if err != nil {
		return "", err
	}

	reader := bzip2.NewReader(bytes.NewReader(source))

	var buf bytes.Buffer

	_, err = io.Copy(&buf, reader)
	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}
