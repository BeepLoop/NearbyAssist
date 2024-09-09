package service

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/store/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserService struct {
	store user.UserStore
}

func NewUserService(store user.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) BaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, "me")
}

func (s *UserService) Login(c echo.Context) error {
	req := new(request.UserLoginPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, req)
}

func (s *UserService) Refresh(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
