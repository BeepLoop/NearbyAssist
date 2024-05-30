package filehandler

import (
	"errors"
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

func FormParser(c echo.Context) ([]*multipart.FileHeader, error) {
	// BUG: hangs when submitted an empty multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	if files, ok := form.File["files"]; !ok {
		return nil, errors.New("No files found in the form")
	} else {
		if len(files) == 0 {
			return nil, err
		}

		return files, nil
	}
}
