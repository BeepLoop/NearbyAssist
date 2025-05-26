package pendingservices

import (
	"context"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/service/activitylog"
	listingreview "nearbyassist/internal/service/listing_review"
	"nearbyassist/internal/service/sse"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/pending_services"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	DEFAULT_LIMIT  = 10
	DEFAULT_OFFSET = 0
)

type handler struct {
	listingReviewService *listingreview.Service
}

func NewHandler(listingReviewService *listingreview.Service) *handler {
	return &handler{
		listingReviewService: listingReviewService,
	}
}

func (h *handler) PendingServicesList(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	params := c.QueryParams()

	limit, _ := strconv.Atoi(params.Get("limit"))
	if limit == 0 {
		limit = DEFAULT_LIMIT
	}

	offset, _ := strconv.Atoi(params.Get("offset"))
	if offset == 0 {
		offset = DEFAULT_OFFSET
	}

	services, err := h.listingReviewService.GetPendingServicesList(limit, offset)
	if err != nil {
		page := pages.PendingServicesList(*admin, make([]dto.PendingServiceListItem, 0), flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.PendingServicesList(*admin, services, flash)
	return page.Render(context.Background(), c.Response().Writer)
}

func (h *handler) PendingServiceDetail(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)
	serviceId := c.Param("serviceId")

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	detail, err := h.listingReviewService.PendingServiceDetail(serviceId)
	if err != nil {
		page := pages.PendingServiceDetail(*admin, dto.PendingService{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.PendingServiceDetail(*admin, *detail, flash)
	return page.Render(context.Background(), c.Response().Writer)
}

func (h *handler) Accept(c echo.Context) error {
	serviceId := c.FormValue("serviceId")
	password := c.FormValue("password")

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	if err := h.listingReviewService.Accept(admin.Id, password, serviceId); err != nil {
		if strings.Contains(err.Error(), listingreview.ERR_UNAUTHORIZED) {
			if err := utils.SetFlashMessage(c, "error", "Unauthorized action"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/pending-services"+serviceId+"?error=unauthorized_action")
			}
		}

		if err := utils.SetFlashMessage(c, "error", "accepting submission failed"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/pending-services"+serviceId+"?error=action_failed")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/pending-services/"+serviceId)
	}

	if err := utils.SetFlashMessage(c, "success", "service submission accepted"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/pending-services?success=submission_accepted")
	}

	activity := activitylog.Input{
		AdminId:    admin.Id,
		Action:     activitylog.ACTION_ACCEPTED_SERVICE,
		TargetType: "other",
		TargetId:   serviceId,
	}
	activitylog.MustGetInstance().CreateWithTarget(activity)
	sse.New().DecreasePendingService()

	return c.Redirect(http.StatusSeeOther, "/admin/pending-services")
}

func (h *handler) Reject(c echo.Context) error {
	serviceId := c.FormValue("serviceId")
	reason := c.FormValue("reason")
	password := c.FormValue("password")

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	if err := h.listingReviewService.Reject(admin.Id, password, serviceId, reason); err != nil {
		if strings.Contains(err.Error(), listingreview.ERR_UNAUTHORIZED) {
			if err := utils.SetFlashMessage(c, "error", "Unauthorized action"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/pending-services/"+serviceId+"?error=unauthorized_action")
			}
		}
		if err := utils.SetFlashMessage(c, "error", "rejecting submission failed"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/pending-services"+serviceId+"?error=action_failed")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/pending-services/"+serviceId)
	}

	if err := utils.SetFlashMessage(c, "success", "service submission rejected"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/pending-services?success=submission_rejected")
	}

	activity := activitylog.Input{
		AdminId:    admin.Id,
		Action:     activitylog.ACTION_REJECTED_SEVICE,
		TargetType: "other",
		TargetId:   serviceId,
	}
	activitylog.MustGetInstance().CreateWithTarget(activity)
	sse.New().DecreasePendingService()

	return c.Redirect(http.StatusSeeOther, "/admin/pending-services")
}
