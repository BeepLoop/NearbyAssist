package fs

import (
	"encoding/hex"
	"errors"
)

type Category string
type FILETYPE string

const (
	STANDARD_JPEG_HEAD_SIGNATURE = "ffd8ffe0"
	EXIF_JPEG_HEAD_SIGNATURE     = "ffd8ffe1"
	SPIFF_JPEG_HEAD_SIGNATURE    = "ffd8ffe8"
	JPEG_TAIL_SIGNATURE          = "ffd9"
	PNG_HEAD_SIGNATURE           = "89504e47"

	FILETYPE_JPEG FILETYPE = "jpeg"
	FILETYPE_PNG  FILETYPE = "png"

	ID_BACK               Category = "id_back"
	ID_FRONT              Category = "id_front"
	FACE                  Category = "face"
	APPLICATION_PROOF_DIR Category = "application_proof"
	SERVICE_PHOTO_DIR     Category = "service_photo"
	SYS_COMPLAINT_DIR     Category = "system_complaint"
	POLICE_CLEARANCE_DIR  Category = "police_clearance"
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

	return FILETYPE_JPEG, nil

	// NOTE: refer to the link for file signatures
	// https://www.garykessler.net/library/file_sigs.html
	// any jpeg format returns jpeg
	if hexForm[:8] == STANDARD_JPEG_HEAD_SIGNATURE || hexForm[:8] == EXIF_JPEG_HEAD_SIGNATURE || hexForm[:8] == SPIFF_JPEG_HEAD_SIGNATURE {
		if hexForm[len(hexForm)-4:] == JPEG_TAIL_SIGNATURE {
			return FILETYPE_JPEG, nil
		}
	}

	if hexForm[:8] == PNG_HEAD_SIGNATURE {
		return FILETYPE_PNG, nil
	}

	return "", errors.New("Unknown filetype")
}
