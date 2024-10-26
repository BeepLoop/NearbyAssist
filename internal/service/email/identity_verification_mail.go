package email

import (
	"bytes"
	"html/template"
	"nearbyassist/internal/response"

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
	var tmpl *template.Template

	switch t {
	case IDENTITY_REQUEST_ACKNOWLEDGMENT:
		if t, err := template.ParseFiles("./views/email/identity_verification_acknowledgment.html"); err != nil {
			return err
		} else {
			tmpl = t
		}
	case IDENTITY_REQUEST_APPROVED:
		if t, err := template.ParseFiles("./views/email/identity_verification_approved.html"); err != nil {
			return err
		} else {
			tmpl = t
		}
	case IDENTITY_REQUEST_REJECTED:
		if t, err := template.ParseFiles("./views/email/identity_verification_rejected.html"); err != nil {
			return err
		} else {
			tmpl = t
		}
	}

	var htmlBytes bytes.Buffer
	if err := tmpl.Execute(&htmlBytes, data); err != nil {
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
