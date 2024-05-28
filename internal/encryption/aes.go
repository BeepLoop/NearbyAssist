package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"nearbyassist/internal/config"
)

type Aes struct {
	key []byte
}

func NewAes(conf *config.Config) *Aes {
	return &Aes{
		key: []byte(conf.EncryptionKey),
	}
}

func (e *Aes) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(ciphertext), nil
}

func (e *Aes) Decrypt(encrypted string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce := encrypted[:nonceSize]
	cipher := encrypted[nonceSize:]

	decrypted, err := gcm.Open(nil, []byte(nonce), []byte(cipher), nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
