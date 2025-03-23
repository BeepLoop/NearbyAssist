package mailer

type MailPayload interface {
	GetContent() string
}
