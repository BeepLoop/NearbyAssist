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

func (e *Aes) Encrypt(source []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	encrypted := gcm.Seal(nonce, nonce, source, nil)

	return encrypted, nil
}

func (e *Aes) Decrypt(source []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce := source[:nonceSize]
	cipher := source[nonceSize:]

	decrypted, err := gcm.Open(nil, nonce, cipher, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func (e *Aes) EncryptString(plaintext string) (string, error) {
	bytes := []byte(plaintext)

	encrypted, err := e.Encrypt(bytes)
	if err != nil {
		return "", err
	}

	return string(encrypted), nil
}

func (e *Aes) DecryptString(encrypted string) (string, error) {
	bytes := []byte(encrypted)

	decrypted, err := e.Decrypt(bytes)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func (e *Aes) EncryptFile(source []byte) ([]byte, error) {
	encrypted, err := e.Encrypt(source)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

func (e *Aes) DecryptFile(source []byte) ([]byte, error) {
	decrypted, err := e.Decrypt(source)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
