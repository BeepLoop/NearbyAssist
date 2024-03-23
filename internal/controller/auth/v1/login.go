package auth

import (
	session_query "nearbyassist/internal/db/query/session"
	user_query "nearbyassist/internal/db/query/user"
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

	_, err := user_query.FindUser(u.Name, u.Email)
	if err != nil {
		id, err := user_query.RegisterUser(*u)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error":   "unable to register user",
				"message": err.Error(),
			})
		}

		u.Id = id
	}

	token, err := utils.GenerateJwt(*u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	err = session_query.NewSession(u.Name, u.Email, token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "could not start a new session",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"userId": u.Id,
		"token":  token,
	})
}
