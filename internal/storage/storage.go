package storage

type Storage interface {
	Initialize() error
	SaveServicePhoto(file []byte, filename string) (string, error)
	SaveApplicationProof(file []byte, filename string) (string, error)
	SaveSystemComplaint(file []byte, filename string) (string, error)
	SaveFrontId(file []byte, filename string) (string, error)
	SaveBackId(file []byte, filename string) (string, error)
	SaveFace(file []byte, filename string) (string, error)
}
