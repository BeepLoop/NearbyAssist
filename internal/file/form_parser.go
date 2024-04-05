package filehandler

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

func FormParser(c echo.Context) ([]*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	files := form.File["files"]
	if len(files) == 0 {
		return nil, err
	}

	return files, nil
}
