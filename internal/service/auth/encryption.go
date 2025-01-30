package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	ENCRYPTION_ERR = "Error occurred while encrypting"
	DECRYPTION_ERR = "Error occurred while decrypting"
)

type Encryption interface {
	EncryptString(text string) (string, error)
	DecryptString(text string) (string, error)
	EncryptFile(source []byte) ([]byte, error)
	DecryptFile(source []byte) ([]byte, error)
	GetKey() ([]byte, error)
}

type MockEncryptor struct{}

func NewMockEncryptor() *MockEncryptor {
	return &MockEncryptor{}
}

func (m *MockEncryptor) EncryptString(text string) (string, error) {
	return text, nil
}

func (m *MockEncryptor) DecryptString(text string) (string, error) {
	return text, nil
}

func (m *MockEncryptor) EncryptFile(source []byte) ([]byte, error) {
	return source, nil
}

func (m *MockEncryptor) DecryptFile(source []byte) ([]byte, error) {
	return source, nil
}

func (m *MockEncryptor) GetKey() ([]byte, error) {
	return nil, nil
}

type AES struct {
	key []byte
}

func NewAES(key []byte) *AES {
	return &AES{key: key}
}

func (e *AES) Encrypt(source []byte) ([]byte, error) {
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

func (e *AES) Decrypt(source []byte) ([]byte, error) {
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

func (e *AES) EncryptString(plaintext string) (string, error) {
	bytes := []byte(plaintext)

	encrypted, err := e.Encrypt(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encrypted), nil
}

func (e *AES) DecryptString(encrypted string) (string, error) {
	bytes, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	decrypted, err := e.Decrypt(bytes)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func (e *AES) EncryptFile(source []byte) ([]byte, error) {
	encrypted, err := e.Encrypt(source)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

func (e *AES) DecryptFile(source []byte) ([]byte, error) {
	decrypted, err := e.Decrypt(source)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
func (e *AES) GetKey() ([]byte, error) {
	if e.key == nil {
		return nil, errors.New("Missing key")
	}

	return e.key, nil
}
