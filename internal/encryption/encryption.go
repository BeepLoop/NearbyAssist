package encryption

type Encryption interface {
	EncryptString(text string) (string, error)
	DecryptString(text string) (string, error)
    EncryptFile(source []byte) ([]byte, error)
    DecryptFile(source []byte) ([]byte, error)
}
