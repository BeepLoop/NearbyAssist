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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := user_query.FindUser(u.Name, u.Email)
	if err != nil {
		id, err := user_query.RegisterUser(*u)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		u.Id = id
	} else {
		u.Id = user.Id
	}

	token, err := utils.GenerateJwt(*u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = session_query.NewSession(u.Name, u.Email, token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"userId": u.Id,
		"token":  token,
	})
}
