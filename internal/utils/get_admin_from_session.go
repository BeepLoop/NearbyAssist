package utils

import (
	"errors"
	"nearbyassist/internal/models"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func GetAdminFromSession(c echo.Context) (*models.AdminModel, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return nil, err
	}

	activeSession := sess.Values["user"]
	admin, ok := activeSession.(models.AdminModel)
	if !ok {
		return nil, errors.New("error binding session to struct")
	}

	return &admin, nil
}
