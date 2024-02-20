package auth

import (
	"nearbyassist/internal/db/query/user"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleLogin(c echo.Context) error {
	u := new(types.User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request data",
		})
	}

	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request data",
		})
	}

	exists, err := query.DoesUserExist(*u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if !exists {
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "account does not exists",
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
