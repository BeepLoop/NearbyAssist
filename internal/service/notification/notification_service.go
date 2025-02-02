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

	NOTIF_CHANNEL_HIGH   = "3b202b57-6bb0-4481-826b-2c9ee405ebf5"
	NOTIF_CHANNEL_URGENT = "133814f9-8cfb-4f14-9b03-7e0d4caa71ba"
)
