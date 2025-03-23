package mailer

import (
	"context"
	"fmt"
	"time"
)

type consoleMailer struct {
	domain string
}

func NewConsoleMailer(domain string) *consoleMailer {
	return &consoleMailer{
		domain: domain,
	}
}

func (m *consoleMailer) Send(parentCtx context.Context, payload MailPayload) error {
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	fmt.Println(payload.GetContent())

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
