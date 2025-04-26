package server

import (
	websocket_handler "nearbyassist/internal/handler/api/websocket"

	"github.com/labstack/echo/v4"
)

func (s *Server) websocketRoute(r *echo.Group) {
	handler := websocket_handler.NewHandler(s.WS, s.JWT)

	r.GET("", handler.Connect)
}
