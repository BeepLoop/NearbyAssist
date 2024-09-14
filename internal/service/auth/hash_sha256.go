package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

type Hash interface {
	Generate(value []byte) (string, error)
}

type Sha256 struct {
	hasher hash.Hash
}

func NewSha256() *Sha256 {
	h := sha256.New()

	return &Sha256{
		hasher: h,
	}
}

func (h *Sha256) Generate(value []byte) (string, error) {
	h.hasher.Reset()

	if _, err := h.hasher.Write(value); err != nil {
		return "", err
	}

	bytes := h.hasher.Sum(nil)
	hash := hex.EncodeToString(bytes)

	return hash, nil
}
