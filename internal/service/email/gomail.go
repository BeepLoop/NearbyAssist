package email

import (
	"errors"

	"gopkg.in/gomail.v2"
)

type GoMail struct {
	from     string
	password string
}

func NewGoMail(email, password string) (*GoMail, error) {
	if email == "" || password == "" {
		return nil, errors.New("Invalid email or password")
	}

	mail := &GoMail{from: email, password: password}

	return mail, nil
}

func (s *GoMail) SendMail(mail Mail) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", s.from)
	msg.SetHeader("To", mail.To...)
	msg.SetHeader("Subject", mail.Subject)

	msg.SetBody("text/html", mail.HtmlMessage)

	msg.AddAlternative("text/plain", mail.PlainMessage)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, s.from, s.password)

	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
