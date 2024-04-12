package storage

import (
	"log"
	"mime/multipart"
)

type DummyStorage struct{}

func newDummyStorage() *DummyStorage {
	return &DummyStorage{}
}

func (s *DummyStorage) Initialize() error {
	log.Print("Initializing dummy storage\n")
	return nil
}

func (s *DummyStorage) SaveServicePhoto(uuid string, file *multipart.FileHeader) (string, error) {
	return "", nil
}

func (s *DummyStorage) SaveApplicationProof(uuid string, file *multipart.FileHeader) (string, error) {
	return "", nil
}
