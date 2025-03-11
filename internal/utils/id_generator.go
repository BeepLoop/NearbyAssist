package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateId() string {
	return gonanoid.Must()
}

func GenerateUserId() string {
	prefix := "user_"
	id := gonanoid.Must()
	return prefix + id
}
