package storage

import (
	"mime/multipart"
)

type Storage interface {
	SaveServicePhoto(uuid string, file *multipart.FileHeader) (string, error)
	SaveApplicationProof(uuid string, file *multipart.FileHeader) (string, error)
}
