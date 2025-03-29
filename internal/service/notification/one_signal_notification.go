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
	instance *oneSignal
)

type oneSignal struct {
	URL    string
	AppID  string
	ApiKey string
}

func NewOneSignal(appId, apiKey string) *oneSignal {
	if instance != nil {
		return instance
	}

	instance = &oneSignal{
		URL:    "https://api.onesignal.com/notifications?c=push",
		AppID:  appId,
		ApiKey: apiKey,
	}

	return instance
}

func MustGetInstance() *oneSignal {
	if instance == nil {
		panic("one signal not instantiated")
	}

	return instance
}

func (n *oneSignal) NewUrgentNotification(userId, heading, content string) error {
	body := MessagePayload{
		AppID: n.AppID,
		Headings: PayloadHeading{
			EN: heading,
		},
		Contents: PayloadContent{
			EN: content,
		},
		IncludeAliases: PayloadIncludeAliases{
			ExternalID: []string{
				userId,
			},
		},
		TargetChannel:    "push",
		AndroidChannelID: NOTIF_CHANNEL_URGENT,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return errors.New(string(NOTIF_PAYLOAD_INIT_ERR))
	}

	return n.shipNotification(bytes.NewReader(payload))
}

func (n *oneSignal) NewMessageNotification(userId string) error {
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
		AndroidChannelID: NOTIF_CHANNEL_URGENT,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return errors.New(string(NOTIF_PAYLOAD_INIT_ERR))
	}

	return n.shipNotification(bytes.NewReader(payload))
}

func (n *oneSignal) shipNotification(payload io.Reader) error {
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
