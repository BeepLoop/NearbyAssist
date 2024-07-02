package handlers

import (
	filehandler "nearbyassist/internal/file"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type complaintHandler struct {
	server *server.Server
}

func NewComplaintHandler(server *server.Server) *complaintHandler {
	return &complaintHandler{
		server: server,
	}
}

func (h *complaintHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Complaints base route",
	})
}

func (h *complaintHandler) HandleSystemComplaintCount(c echo.Context) error {
	count, err := h.server.DB.CountSystemComplaint()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"count": count,
	})
}

func (h *complaintHandler) HandleGetSystemComplaint(c echo.Context) error {
	complaints, err := h.server.DB.FindAllSystemComplaints()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"complaints": complaints,
	})
}

func (h *complaintHandler) HandleGetSystemComplaintById(c echo.Context) error {
	complaintId := c.Param("complaintId")
	id, err := strconv.Atoi(complaintId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "complaint ID must be a number")
	}

	complaint, err := h.server.DB.FindSystemComplaintById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	images, err := h.server.DB.FindSystemComplaintImagesByComplaintId(complaint.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"complaint": complaint,
		"images":    images,
	})
}

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

	imageUrl := make([]string, 0)
	for _, file := range files {
		handler := filehandler.NewFileHandler(h.server.Encrypt)
		url, err := handler.SavePhoto(file, h.server.Storage.SaveSystemComplaint)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		imageUrl = append(imageUrl, url)
	}

	for _, url := range imageUrl {
		model := &models.SystemComplaintImageModel{
			ComplaintId: complaintId,
			Url:         url,
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
	// TODO: implement filing a complaint for a vendor
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Vendor complaint route",
	})
}
