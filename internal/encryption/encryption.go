package encryption

type Encryption interface {
	Encrypt(text string) (string, error)
	Decrypt(text string) (string, error)
}
