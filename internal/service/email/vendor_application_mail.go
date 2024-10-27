package email

import (
	"bytes"
	"context"
	"nearbyassist/internal/response"
	"nearbyassist/views/email"

	"github.com/a-h/templ"
	"github.com/labstack/gommon/log"
)

type VendorApplicationMailType string

const (
	VENDOR_APPLICATION_ACKNOWLEDGMENT VendorApplicationMailType = "Vendor Application Request Acknowledgment"
	VENDOR_APPLICATION_APPROVED       VendorApplicationMailType = "Vendor Application Request Approved"
	VENDOR_APPLICATION_REJECTED       VendorApplicationMailType = "Vendor Application Request Rejected"
)

type vendorApplicationMail struct {
	mail   Mail
	mailer MailService
}

func VendorApplicationMail(mailer MailService) *vendorApplicationMail {
	return &vendorApplicationMail{
		mail:   Mail{},
		mailer: mailer,
	}
}

func (m *vendorApplicationMail) SetBody(t VendorApplicationMailType, data response.BasicEmailPayload) error {
	var view templ.Component

	switch t {
	case VENDOR_APPLICATION_ACKNOWLEDGMENT:
		view = email.VendorApplicationAcknowledgment(data)
	case VENDOR_APPLICATION_APPROVED:
		view = email.VendorApplicationApproved(data)
	case VENDOR_APPLICATION_REJECTED:
		view = email.VendorApplicationRejected(data)
	}

	var htmlBytes bytes.Buffer
	if err := view.Render(context.Background(), &htmlBytes); err != nil {
		return err
	}

	m.mail.HtmlBody = htmlBytes.String()

	return nil
}

func (m *vendorApplicationMail) SetSubject(subject VendorApplicationMailType) *vendorApplicationMail {
	m.mail.Subject = string(subject)
	return m
}

func (m *vendorApplicationMail) To(to []string) *vendorApplicationMail {
	m.mail.To = to
	return m
}

func (m *vendorApplicationMail) Send() error {
	err := m.mailer.SendMail(m.mail)
	if err != nil {
		log.Info("Error sending vendor application mail")
		return err
	}

	return nil
}
