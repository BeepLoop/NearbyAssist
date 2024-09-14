package fs

import (
	"encoding/hex"
	"errors"
)

type Category string
type FILETYPE string

const (
	JPEG_HEAD_SIGNATURE = "ffd8"
	JPEG_TAIL_SIGNATURE = "ffd9"
	PNG_HEAD_SIGNATURE  = "89504e47"

	FILETYPE_JPEG FILETYPE = "jpeg"
	FILETYPE_PNG  FILETYPE = "png"

	ID_BACK               Category = "id_back"
	ID_FRONT              Category = "id_front"
	FACE                  Category = "face"
	APPLICATION_PROOF_DIR Category = "application_proof"
	SERVICE_PHOTO_DIR     Category = "service_photo"
	SYS_COMPLAINT_DIR     Category = "system_complaint"
)

type File struct {
	Data     []byte
	Category Category
}

type FileStorage interface {
	SaveFile(file File) (string, error)
	GetFile(path string) ([]byte, error)
}

func GetFiletype(file []byte) (FILETYPE, error) {
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
