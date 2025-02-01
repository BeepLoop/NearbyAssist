package notification_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	OneSignalInstance *OneSignalNotification
)

type OneSignalNotification struct {
	URL    string
	AppID  string
	ApiKey string
}

func NewOneSignal(appId, apiKey string) *OneSignalNotification {
	if OneSignalInstance != nil {
		return OneSignalInstance
	}

	OneSignalInstance = &OneSignalNotification{
		URL:    "https://api.onesignal.com/notifications?c=push",
		AppID:  appId,
		ApiKey: apiKey,
	}

	return OneSignalInstance
}

func (n *OneSignalNotification) NewMessageNotification(userId string, notifyType NotificationType) error {
	body := MessagePayload{
		AppID: n.AppID,
		Headings: PayloadHeading{
			EN: "New Message",
		},
		Contents: PayloadContent{
			EN: "1 new unread message",
		},
		IncludeAliases: PayloadIncludeAliases{
			ExternalID: []string{
				userId,
			},
		},
		TargetChannel:    "push",
		AndroidChannelID: "133814f9-8cfb-4f14-9b03-7e0d4caa71ba",
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return errors.New(string(NOTIF_PAYLOAD_INIT_ERR))
	}

	return n.shipNotification(bytes.NewReader(payload))
}

func (n *OneSignalNotification) shipNotification(payload io.Reader) error {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.URL, payload)
	if err != nil {
		return errors.New(string(NOTIF_REQ_INIT_ERR))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Key "+n.ApiKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New(string(NOTIF_HTTP_CLIENT_ERR))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(NOTIF_PAYLOAD_ERR))
	}

	return nil
}

func (n *OneSignalNotification) payloadFactory(userId string, notifType NotificationType) (io.Reader, error) {
	switch notifType {
	case NOTIF_TYPE_NEW_MESSAGE:
		return n.newMessagePayload(userId)
	default:
		return nil, errors.New(string(NOTIF_TYPE_INVALID))
	}
}

func (n *OneSignalNotification) newMessagePayload(userId string) (io.Reader, error) {
	body := MessagePayload{
		AppID: n.AppID,
		Headings: PayloadHeading{
			EN: "New Message",
		},
		Contents: PayloadContent{
			EN: "1 new unread message",
		},
		IncludeAliases: PayloadIncludeAliases{
			ExternalID: []string{
				userId,
			},
		},
		TargetChannel: "push",
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(payload), nil
}
