package mailer

import "context"

type Mailer interface {
	Send(context.Context, MailPayload) error
}
