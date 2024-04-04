package upload

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func VendorApplicationProof(c echo.Context) error {
	// TODO: get uplaod params

	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	files := form.File["files"]
	if len(files) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "no files attached")
	}

	for _, file := range files {
		log.Println(file.Filename)
		// TODO: handle filesaver upload call
	}
	return c.JSON(http.StatusCreated, "Vendor application proof uploaded successfully")
}
