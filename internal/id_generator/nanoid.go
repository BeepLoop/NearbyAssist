package id_generator

import gonanoid "github.com/matoous/go-nanoid/v2"

type NanoId struct{}

func NewNanoIdGenerator() *NanoId {
	return &NanoId{}
}

func (n *NanoId) Generate() (string, error) {
	id, err := gonanoid.New()
	if err != nil {
		return "", err
	}

	return id, nil
}
