package email

import (
	"bytes"
	"html/template"
	"nearbyassist/internal/response"

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
	var tmpl *template.Template

	switch t {
	case VENDOR_APPLICATION_ACKNOWLEDGMENT:
		if t, err := template.ParseFiles("./vendor_application_acknowledgment.html"); err != nil {
			return err
		} else {
			tmpl = t
		}
	case VENDOR_APPLICATION_APPROVED:
		if t, err := template.ParseFiles("./vendor_application_approved.html"); err != nil {
			return err
		} else {
			tmpl = t
		}
	case VENDOR_APPLICATION_REJECTED:
		if t, err := template.ParseFiles("./vendor_application_rejected.html"); err != nil {
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
