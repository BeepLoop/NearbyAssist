package email

import (
	"bytes"
	"context"
	"nearbyassist/internal/response"
	"nearbyassist/views/email"

	"github.com/labstack/gommon/log"
)

type transactionSummaryMail struct {
	mail   Mail
	mailer MailService
}

func TransactionSummaryMail(mailer MailService) *transactionSummaryMail {
	return &transactionSummaryMail{
		mail: Mail{
			Subject: "Initiated Transaction Summary",
		},
		mailer: mailer,
	}
}

func (m *transactionSummaryMail) SetBody(data response.TransactionSummary) error {
	view := email.TransactionSummary(data)

	var htmlBytes bytes.Buffer
	if err := view.Render(context.Background(), &htmlBytes); err != nil {
		return err
	}

	m.mail.HtmlBody = htmlBytes.String()

	return nil
}

func (m *transactionSummaryMail) To(to []string) *transactionSummaryMail {
	m.mail.To = to
	return m
}

func (m *transactionSummaryMail) Send() error {
	err := m.mailer.SendMail(m.mail)
	if err != nil {
		log.Info("Error sending vendor application mail")
		return err
	}

	return nil
}
