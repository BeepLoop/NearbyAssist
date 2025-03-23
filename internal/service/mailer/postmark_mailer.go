package mailer

import (
	"context"
	"fmt"
	"time"

	"github.com/mrz1836/postmark"
)

type postmarkMailer struct {
	serverToken  string
	accountToken string
}

func NewPostmarkMailer(serverToken, accountToken string) *postmarkMailer {
	return &postmarkMailer{
		serverToken:  serverToken,
		accountToken: accountToken,
	}
}

func (m *postmarkMailer) Send(parentCtx context.Context, payload MailPayload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	email := postmark.Email{
		From:       "no-reply@johnloydm.site",
		To:         payload.GetRecipient(),
		Subject:    payload.GetSubject(),
		HTMLBody:   payload.GetHTML(),
		TextBody:   "Email from NearbyAssist team.",
		Tag:        "invitation",
		TrackOpens: true,
	}

	client := postmark.NewClient(m.serverToken, m.accountToken)

	_, err := client.SendEmail(ctx, email)
	if err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
