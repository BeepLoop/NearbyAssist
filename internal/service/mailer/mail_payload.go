package mailer

type MailPayload interface {
	GetHTML() string
	GetRecipient() string
	GetSubject() string
}
