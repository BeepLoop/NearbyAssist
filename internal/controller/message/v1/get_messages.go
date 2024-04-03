package message

import (
	"nearbyassist/internal/db/query/message"
	"nearbyassist/internal/types"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetMessages(c echo.Context) error {
	fromId := c.QueryParam("from")
	toId := c.QueryParam("to")
	if fromId == "" || toId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "required params field missing")
	}

	from, err := strconv.Atoi(fromId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "from ID must be a number")
	}

	to, err := strconv.Atoi(toId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "to ID must be a number")
	}

	msgParams := types.Message{
		Sender:   from,
		Receiver: to,
	}

	messages, err := message_query.GetMessages(msgParams)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, messages)
}
