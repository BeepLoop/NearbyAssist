package websocket_handler

import (
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
	"net/http"

	gorilla_ws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type handler struct {
	ws  websocket.Socket
	jwt core.Authenticator
}

func NewHandler(ws websocket.Socket, jwt core.Authenticator) *handler {
	return &handler{
		ws:  ws,
		jwt: jwt,
	}
}

func (h *handler) Connect(c echo.Context) error {
	conn, err := h.ws.Upgrade(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
			Message: "Error upgrading connection",
			Error:   err.Error(),
		})
	}
	defer conn.Close()

	// NOTE: Use bearer token when I successfully solved the problem of passing
	// JWT from client Authorization header that works even on reconnect and JWT
	// updates.
	token := c.QueryParam("token")
	userId, err := utils.GetUserIdFromToken(token, h.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
			Message: "Invalid token",
			Error:   err.Error(),
		})
	}

	h.ws.RegisterClient(userId, conn)
	fmt.Println("user connected: ", userId)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if gorilla_ws.IsCloseError(err, gorilla_ws.CloseNormalClosure, gorilla_ws.CloseGoingAway, gorilla_ws.CloseAbnormalClosure) {
				h.ws.RemoveClient(userId)
				fmt.Println("User disconnected: ", userId)

				return nil
			}

			fmt.Println("read error: ", err.Error())
			return err
		}

		fmt.Println("received: ", string(msg))
	}
}
