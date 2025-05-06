package ssehandler

import (
	"fmt"
	"nearbyassist/internal/service/sse"
	"time"

	"github.com/labstack/echo/v4"
)

type handler struct {
	clients []string
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Listen(c echo.Context) error {
	fmt.Println("sse client connected, ip: %v", c.RealIP())

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request().Context().Done():
			fmt.Println("sse client disconnected, ip: %v", c.RealIP())
			return nil
		case <-ticker.C:
			b, err := sse.New().GetMarshalled()
			if err != nil {
				fmt.Println("error marshalling sse data: ", err.Error())
				return err
			}

			data := fmt.Sprintf("data: %s\n\n", b)

			c.Response().Writer.Write([]byte(data))
			c.Response().Flush()
		}
	}
}
