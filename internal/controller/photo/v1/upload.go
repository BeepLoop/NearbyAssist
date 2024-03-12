package photo

import (
	"fmt"
	photo_query "nearbyassist/internal/db/query/photo"
	"nearbyassist/internal/types"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadImage(c echo.Context) error {
	vendorId, serviceId, err := utils.GetUploadParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

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
		// Save file to local storage
		filename, err := utils.FileSaver(file, vendorId, serviceId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		// Save file location to database
		fileData := types.UploadData{
			VendorId:  vendorId,
			ServiceId: serviceId,
			ImageUrl:  fmt.Sprintf("/resource/%s", filename),
		}
		err = photo_query.UploadPhoto(fileData)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "upload successful",
	})
}
