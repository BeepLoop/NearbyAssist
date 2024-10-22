package email

type Mail struct {
	To              []string
	Subject         string
	HtmlBody        string
	AlternativeBody string
}

type MailService interface {
	SendMail(mail Mail) error
}
