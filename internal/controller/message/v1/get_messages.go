package message

import (
	query "nearbyassist/internal/db/query/message"
	"nearbyassist/internal/types"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetMessages(c echo.Context) error {
	fromId := c.QueryParam("from")
	toId := c.QueryParam("to")
	if fromId == "" || toId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	from, err := strconv.Atoi(fromId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	to, err := strconv.Atoi(toId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	msgParams := types.Message{
		Sender:   from,
		Reciever: to,
	}

	messages, err := query.GetMessages(msgParams)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, messages)
}
