package storage

import (
	"log"
)

type DummyStorage struct{}

func newDummyStorage() *DummyStorage {
	return &DummyStorage{}
}

func (s *DummyStorage) Initialize() error {
	log.Print("Initializing dummy storage\n")
	return nil
}

func (s *DummyStorage) SaveServicePhoto(file []byte, filename string) (string, error) {
	return "", nil
}

func (s *DummyStorage) SaveApplicationProof(file []byte, filename string) (string, error) {
	return "", nil
}

func (s *DummyStorage) SaveSystemComplaint(file []byte, filename string) (string, error) {
	return "", nil
}
