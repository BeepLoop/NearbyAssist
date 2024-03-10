package photo

import (
	"fmt"
	"io"
	photo_query "nearbyassist/internal/db/query/photo"
	"nearbyassist/internal/types"
	"nearbyassist/internal/utils"
	"net/http"
	"os"
	"strings"
	"time"

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
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		defer src.Close()

		timestamp := time.Now().Format("2006-01-02_15:04:05")
		mimeType := strings.Split(file.Header["Content-Type"][0], "/")[1]
		distFilename := fmt.Sprintf("%d_%d_%s.%s", vendorId, serviceId, timestamp, mimeType)

		// create the file in the server
		dist, err := os.Create("store/" + distFilename)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		defer dist.Close()

		// copy the uploaded file to the opened file
		_, err = io.Copy(dist, src)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		// Save file location to database
		fileData := types.UploadData{
			VendorId:  vendorId,
			ServiceId: serviceId,
			ImageUrl:  fmt.Sprintf("/resource/%s", distFilename),
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
