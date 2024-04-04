package filehandler

type FileSaver interface {
	SavePhoto(uuid string) (string, error)
	Upload() error
}
