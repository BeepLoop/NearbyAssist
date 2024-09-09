package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(value []byte) (string, error) {
	h := sha256.New()
	h.Reset()

	if _, err := h.Write(value); err != nil {
		return "", err
	}

	bytes := h.Sum(nil)
	hash := hex.EncodeToString(bytes)

	return hash, nil
}
