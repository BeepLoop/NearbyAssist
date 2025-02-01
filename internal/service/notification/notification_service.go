package notification_service

type NotificationType string
type NotificationError string

const (
	NOTIF_TYPE_NEW_MESSAGE NotificationType = "New Message"

	NOTIF_TYPE_INVALID     NotificationError = "Invalid notification type"
	NOTIF_REQ_INIT_ERR     NotificationError = "Error initializing http request"
	NOTIF_PAYLOAD_INIT_ERR NotificationError = "Error initializing payload"
	NOTIF_PAYLOAD_ERR      NotificationError = "Payload error"
	NOTIF_HTTP_CLIENT_ERR  NotificationError = "Error sending notification"
)

type NotificationService interface {
	Notify(userId string, notifType NotificationType) error
}
