package filehandler

type FileSaver interface {
	SavePhoto() (string, error)
}
