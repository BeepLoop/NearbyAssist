package email

import (
	"bytes"
	"html/template"
	"nearbyassist/internal/response"

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
	tmpl, err := template.ParseFiles("./views/email/transaction_summary.html")
	if err != nil {
		return err
	}

	var htmlBytes bytes.Buffer
	if err := tmpl.Execute(&htmlBytes, data); err != nil {
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
