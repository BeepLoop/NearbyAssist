package email

type Mail struct {
	To           []string
	Subject      string
	PlainMessage string
	HtmlMessage  string
}

type Email interface {
	SendMail(mail Mail) error
}
