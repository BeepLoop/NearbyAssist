package mailer

import (
	"context"
	"fmt"
	"time"
)

type consoleMailer struct{}

func NewConsoleMailer() *consoleMailer {
	return &consoleMailer{}
}

func (m *consoleMailer) Send(parentCtx context.Context, payload MailPayload) error {
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	fmt.Println(payload.GetHTML())

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
