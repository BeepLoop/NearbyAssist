package ssehandler

import (
	"fmt"
	"nearbyassist/internal/service/sse"
	"net/http"
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
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	f, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case <-ticker.C:
			b, err := sse.New().GetMarshalled()
			if err != nil {
				fmt.Println("error marshalling sse data: ", err.Error())
				return err
			}

			fmt.Fprintf(c.Response().Writer, "data: %s\n\n", b)
			f.Flush()
		}
	}
}
