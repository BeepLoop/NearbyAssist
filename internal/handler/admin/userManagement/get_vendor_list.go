package userManagement

import (
	"context"
	"fmt"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/user_management/seller"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) GetVendorList(c echo.Context) error {
	params := c.QueryParams()

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	limit, _ := strconv.Atoi(params.Get("limit"))
	if limit == 0 {
		limit = DEFAULT_LIMIT
	}

	offset, _ := strconv.Atoi(params.Get("offset"))
	if offset == 0 {
		offset = DEFAULT_OFFSET
	}

	results := make([]dto.Vendor, 0)

	// If 'query' exists, perform search
	if params.Has("query") && params.Get("query") != "" {
		if strings.HasPrefix(params.Get("query"), "email:") {
			query := params.Get("query")[len("email:"):]

			user, err := h.vendorService.GetAllByEmail(query)
			if err != nil {
				fmt.Println("error get vendors by email: ", err.Error())
				page := pages.VendorList(*admin, make([]dto.Vendor, 0))
				return page.Render(context.Background(), c.Response().Writer)
			}
			results = append(results, *user)
		} else if strings.HasPrefix(params.Get("query"), "title:") {
			query := params.Get("query")[len("title:"):]

			users, err := h.vendorService.GetAllByExpertise(query, limit, offset)
			if err != nil {
				fmt.Println("error get vendors by expertise: ", err.Error())
				page := pages.VendorList(*admin, make([]dto.Vendor, 0))
				return page.Render(context.Background(), c.Response().Writer)
			}
			results = users
		} else if strings.HasPrefix(params.Get("query"), "id:") {
			query := params.Get("query")[len("id:"):]

			user, err := h.vendorService.FindByIdDTO(query)
			if err != nil {
				fmt.Println("error get vendors by expertise: ", err.Error())
				page := pages.VendorList(*admin, make([]dto.Vendor, 0))
				return page.Render(context.Background(), c.Response().Writer)
			}

			results = append(results, *user)
		} else {
			// Fallback search by email
			query := params.Get("query")

			user, err := h.vendorService.GetAllByEmail(query)
			if err != nil {
				fmt.Println("error get vendors by email: ", err.Error())
				page := pages.VendorList(*admin, make([]dto.Vendor, 0))
				return page.Render(context.Background(), c.Response().Writer)
			}
			results = append(results, *user)
		}
	} else {
		accounts, err := h.vendorService.GetAll(limit, offset)
		if err != nil {
			fmt.Println("error get all vendors: ", err.Error())
			page := pages.VendorList(*admin, make([]dto.Vendor, 0))
			return page.Render(context.Background(), c.Response().Writer)
		}
		results = accounts
	}

	page := pages.VendorList(*admin, results)
	return page.Render(context.Background(), c.Response().Writer)
}
