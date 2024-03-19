package auth

import (
	// "nearbyassist/internal/db/query/user"
	session_query "nearbyassist/internal/db/query/session"
	"nearbyassist/internal/types"
	"nearbyassist/internal/utils"
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

	token, err := utils.GenerateJwt(*u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	err = session_query.NewSession(u.Name, token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "active session already exists",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"token":  token,
	})
}
