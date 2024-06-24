package handlers

import (
	filehandler "nearbyassist/internal/file"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type complaintHandler struct {
	server *server.Server
}

func NewComplaintServer(server *server.Server) *complaintHandler {
	return &complaintHandler{
		server: server,
	}
}

func (h *complaintHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Complaints base route",
	})
}

func (h *complaintHandler) HandleCount(c echo.Context) error {
	count, err := h.server.DB.CountComplaint()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"complaintsCount": count,
	})
}

// func (h *complaintHandler) HandleNewComplaint(c echo.Context) error {
// 	req := &request.NewComplaint{}
// 	if err := c.Bind(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}
//
// 	if err := c.Validate(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}
//
// 	complaintId, err := h.server.DB.FileComplaint(req)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 	}
//
// 	return c.JSON(http.StatusCreated, utils.Mapper{
// 		"message":     "complaint created successfully",
// 		"complaintId": complaintId,
// 	})
// }

func (h *complaintHandler) HandleSystemComplaint(c echo.Context) error {
	title := c.FormValue("title")
	detail := c.FormValue("detail")

	files, err := filehandler.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	req := &request.SystemComplaint{
		Title:  title,
		Detail: detail,
	}

	complaintId, err := h.server.DB.FileSystemComplaint(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	supportingImages := make([]string, 0)
	for _, file := range files {
		handler := filehandler.NewFileHandler(h.server.Storage)
		filename, err := handler.SavePhoto(file, h.server.Storage.SaveSystemComplaint)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		supportingImages = append(supportingImages, filename)
	}

	for _, filename := range supportingImages {
		model := &models.SystemComplaintImageModel{
			ComplaintId: complaintId,
			Url:         filename,
		}

		_, err := h.server.DB.NewSystemComplaintImage(model)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"message": "System complaint created successfully",
	})
}

func (h *complaintHandler) HandleVendorComplaint(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Vendor complaint route",
	})
}
