package mailer

import "nearbyassist/internal/service/email"

type Mailer struct {
	mailChan chan email.Mail
	mail     email.Email
}

func New(mail email.Email) *Mailer {
	return &Mailer{mailChan: make(chan email.Mail)}
}

func (m *Mailer) Start() {
	go func() {
		for {
			select {
			case mail := <-m.mailChan:
				m.mail.SendMail(mail)
			}
		}
	}()
}

func (m *Mailer) Add(mail email.Mail) {
	m.mailChan <- mail
}

func (m *Mailer) Stop() {
	close(m.mailChan)
}
