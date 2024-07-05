package hash

const (
	HASH_ERROR = "Error occurred while processing your data"
)

type Hash interface {
	Hash(value []byte) (string, error)
}
