package utils

import (
	"io"
	"mime/multipart"
)

func FileToBytes(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
