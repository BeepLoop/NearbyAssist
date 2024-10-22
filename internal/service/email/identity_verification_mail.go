package email

import "github.com/labstack/gommon/log"

type identityVerificationMail struct {
	mail   Mail
	mailer MailService
}

func IdentityVerificationMail(mailer MailService) *identityVerificationMail {
	return &identityVerificationMail{
		mail: Mail{
			Subject: "Identity Verification Request Acknowledgement",
			HtmlBody: `
                <html>
                    <body>
                        <h1>Identity Verification Request Acknowledgement</h1>
                        <h3>Thank you for using NearbyAssist</h3>
                        <p>We have received your request for identity verification, we are processing your request and we will get back to you.</p>
                    </body>
                </html>
            `,
			AlternativeBody: "Identity Verification Request Acknowledgement\n\nThank you for using NearbyAssist. We have received your request for identity verification, we are processing your request and we will get back to you.",
		},
		mailer: mailer,
	}
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
