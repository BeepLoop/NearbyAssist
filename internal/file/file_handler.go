package filehandler

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"nearbyassist/internal/encryption"

	"github.com/google/uuid"
)

type FILETYPE string

const (
	JPEG_HEAD_SIGNATURE = "ffd8"
	JPEG_TAIL_SIGNATURE = "ffd9"
	PNG_HEAD_SIGNATURE  = "89504e47"

	FILETYPE_JPEG FILETYPE = "jpeg"
	FILETYPE_PNG  FILETYPE = "png"
)

type FileHandler struct {
	encryptor encryption.Encryption
}

func NewFileHandler(encryptor encryption.Encryption) *FileHandler {
	return &FileHandler{
		encryptor: encryptor,
	}
}

func (f *FileHandler) SavePhoto(file *multipart.FileHeader, saveFunc func(src []byte, filename string) (string, error)) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}

	bytes, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	encrypted, err := f.encryptor.EncryptFile(bytes)
	if err != nil {
		return "", err
	}

	uuid := uuid.New().String()
	extension, err := f.GetFileExtension(bytes)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.%s", uuid, extension)

	if _, err := saveFunc(encrypted, filename); err != nil {
		return "", err
	}

	return filename, nil
}

func (f *FileHandler) GetFileExtension(file []byte) (FILETYPE, error) {
	hexForm := hex.EncodeToString(file)

	// NOTE: refer to the link for file signatures
	// https://www.garykessler.net/library/file_sigs.html
	if hexForm[:4] == JPEG_HEAD_SIGNATURE && hexForm[len(hexForm)-4:] == JPEG_TAIL_SIGNATURE {
		return FILETYPE_JPEG, nil
	}

	if hexForm[:8] == PNG_HEAD_SIGNATURE {
		return FILETYPE_PNG, nil
	}

	return "", errors.New("Unknown filetype")
}
