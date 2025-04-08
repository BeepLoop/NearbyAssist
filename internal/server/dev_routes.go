package server

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"

	"github.com/labstack/echo/v4"
)

func (s *Server) DevRoutes(r *echo.Group) {

	pingRoute := r.Group("/ping")
	{
		pingRoute.GET("/:userId", func(c echo.Context) error {
			userId := c.Param("userId")

			notification := &models.NotificationModel{
				Model:     models.Model{Id: utils.GenerateId()},
				Recipient: userId,
				Type:      "success",
				Title:     "Test notification",
				Content:   "This is a test notification",
			}

			notifEvent := &websocket.EventModel{
				ReceiverId: userId,
				Type:       websocket.EVT_NOTIF,
				Payload:    notification,
			}

			s.WS.Send(notifEvent)

			return nil
		})
	}
}
