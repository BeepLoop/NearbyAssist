package encryption

const (
	ENCRYPTION_ERR = "Error occurred while encrypting"
	DECRYPTION_ERR = "Error occurred while decrypting"
)

type Encryption interface {
	EncryptString(text string) (string, error)
	DecryptString(text string) (string, error)
	EncryptFile(source []byte) ([]byte, error)
	DecryptFile(source []byte) ([]byte, error)
}
