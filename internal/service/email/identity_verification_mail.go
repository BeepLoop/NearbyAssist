package email

import (
	"bytes"
	"context"
	"nearbyassist/internal/response"
	"nearbyassist/views/email"

	"github.com/a-h/templ"
	"github.com/labstack/gommon/log"
)

type IdentityVerificationMailType string

const (
	IDENTITY_REQUEST_ACKNOWLEDGMENT IdentityVerificationMailType = "Identity Verification Request Acknowledgment"
	IDENTITY_REQUEST_APPROVED       IdentityVerificationMailType = "Identity Verification Request Approved"
	IDENTITY_REQUEST_REJECTED       IdentityVerificationMailType = "Identity Verification Request Rejected"
)

type identityVerificationMail struct {
	mail   Mail
	mailer MailService
}

func IdentityVerificationMail(mailer MailService) *identityVerificationMail {
	return &identityVerificationMail{
		mail:   Mail{},
		mailer: mailer,
	}
}

func (m *identityVerificationMail) SetBody(t IdentityVerificationMailType, data response.BasicEmailPayload) error {
	var view templ.Component

	switch t {
	case IDENTITY_REQUEST_ACKNOWLEDGMENT:
		view = email.IdentityVerificationAcknowledgment(data)
	case IDENTITY_REQUEST_APPROVED:
		view = email.IdentityVerificationApproved(data)
	case IDENTITY_REQUEST_REJECTED:
		view = email.IdentityVerificationRejected(data)
	}

	var htmlBytes bytes.Buffer
	if err := view.Render(context.Background(), &htmlBytes); err != nil {
		return err
	}

	m.mail.HtmlBody = htmlBytes.String()

	return nil
}

func (m *identityVerificationMail) SetSubject(subject IdentityVerificationMailType) *identityVerificationMail {
	m.mail.Subject = string(subject)
	return m
}

func (m *identityVerificationMail) To(to []string) *identityVerificationMail {
	m.mail.To = to
	return m
}

func (m *identityVerificationMail) Send() error {
	err := m.mailer.SendMail(m.mail)
	if err != nil {
		log.Info("Error sending identity verification mail")
		return err
	}

	return nil
}
