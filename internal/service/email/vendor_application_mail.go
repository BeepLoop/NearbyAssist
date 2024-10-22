package email

import "github.com/labstack/gommon/log"

type vendorApplicationMail struct {
	mail   Mail
	mailer MailService
}

func VendorApplicationMail(mailer MailService) *vendorApplicationMail {
	return &vendorApplicationMail{
		mail: Mail{
			Subject: "Vendor Application Request Acknowledgement",
			HtmlBody: `
                <html>
                    <body>
                        <h1>Vendor Application Request Acknowledgement</h1>
                        <h3>Thank you for using NearbyAssist</h3>
                        <p>We have received your vendor application request, we are processing your request and we will get back to you.</p>
                    </body>
                </html>
            `,
			AlternativeBody: "Identity Verification Request Acknowledgement\n\nThank you for using NearbyAssist. We have received your vendor application request, we are processing your request and we will get back to you.",
		},
		mailer: mailer,
	}
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
