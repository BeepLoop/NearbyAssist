package auth

import gonanoid "github.com/matoous/go-nanoid/v2"

const (
	NANO_ID_ERR = "Error generating nano id"
)

func GenerateNanoId() (string, error) {
	if id, err := gonanoid.New(); err != nil {
		return "", err
	} else {
		return id, nil
	}
}
