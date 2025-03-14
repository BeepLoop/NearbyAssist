package resource_service

import "nearbyassist/internal/service/fs"

type Service struct {
	fs fs.FileStorage
}

func NewService(fs fs.FileStorage) *Service {
	return &Service{
		fs: fs,
	}
}

func (s *Service) GetFile(path string) ([]byte, error) {
	// Retrieve the file
	b, err := s.fs.GetFile(path)
	if err != nil {
		return nil, err
	}

	return b, nil
}
