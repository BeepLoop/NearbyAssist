package photo

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func UploadImage(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "no files attached",
		})
	}

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		defer src.Close()

		// create the file in the server
		dist, err := os.Create(file.Filename)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		defer dist.Close()

		// copy the uploaded file to the opened file
		if _, err = io.Copy(dist, src); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "upload successful",
	})
}
