package auth

import (
	"nearbyassist/internal/db/query/user"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleRegister(c echo.Context) error {
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

	resultPtr, err := query.GetUser(*u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if resultPtr == nil {
		err = query.RegisterUser(*u)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
	}

	return c.JSON(http.StatusCreated, u)
}
