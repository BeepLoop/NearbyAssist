package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

type sha struct {
	hasher hash.Hash
}

func NewSha() *sha {
	hasher := sha256.New()

	return &sha{hasher: hasher}
}

func (h *sha) Hash(value []byte) (string, error) {
	if _, err := h.hasher.Write(value); err != nil {
		return "", err
	}

	bytes := h.hasher.Sum(nil)
	hash := hex.EncodeToString(bytes)

	return hash, nil
}
