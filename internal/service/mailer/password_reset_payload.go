package mailer

import (
	"bytes"
	"text/template"
)

type PasswordResetPayload struct {
	Username  string
	Password  string
	Email     string
	LoginLink string
}

func (p *PasswordResetPayload) GetHTML() string {
	emailTemplate := `
        <!DOCTYPE html>
        <html>
        <head>
          <meta charset="UTF-8">
          <meta name="viewport" content="width=device-width, initial-scale=1.0">
          <title>Password Reset Successful</title>
        </head>
        <body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4; text-align: center;">
          <table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0">
            <tr>
              <td align="center">
                <table role="presentation" width="100%" max-width="600px" cellspacing="0" cellpadding="0" border="0" style="background-color: #ffffff; border-radius: 10px; padding: 20px; margin: 20px;">
                  <!-- Header -->
                  <tr>
                    <td align="center" style="padding: 20px 0;">
                      <h1 style="color: #2E7D32;">NearbyAssist</h1>
                    </td>
                  </tr>
                  <!-- Notification Message -->
                  <tr>
                    <td align="center">
                      <h2 style="color: #333;">Password Reset Successful 🔒</h2>
                      <p style="font-size: 18px; color: #666;">Hello <strong>{{ .Username }}</strong>,</p>
                      <p style="font-size: 16px; color: #666;">Your password has been successfully reset. Below are your new login credentials:</p>
                    </td>
                  </tr>
                  <!-- Login Details -->
                    <tr>
                        <td align="center" style="padding: 10px;">
                            <table role="presentation" width="80%" cellspacing="0" cellpadding="10" border="0"
                                style="background-color: #E8F5E9; border-radius: 5px;">
                                <tr>
                                    <td style="color: #2E7D32; font-size: 16px; text-align: left;"><strong>👤
                                            Username:</strong></td>
                                    <td style="color: #333; font-size: 16px; text-align: right;">
                                        <strong>{{ .Username }}</strong>
                                    </td>
                                </tr>
                                <tr>
                                    <td style="color: #2E7D32; font-size: 16px; text-align: left;"><strong>🔑
                                            Password:</strong></td>
                                    <td style="color: #333; font-size: 16px; text-align: right;">
                                        <strong>{{ .Password }}</strong>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                  <!-- Note -->
                  <tr>
                    <td align="center" style="padding: 20px;">
                      <p style="font-size: 16px; color: #d32f2f;"><strong>⚠️ Important:</strong> You are required to change your password upon first login for security purposes.</p>
                    </td>
                  </tr>
                  <!-- Call to Action -->
                  <tr>
                    <td align="center" style="padding: 20px 0;">
                      <a href="{{ .LoginLink }}" style="background-color: #2E7D32; color: #ffffff; text-decoration: none; padding: 12px 24px; border-radius: 5px; font-size: 16px; display: inline-block;">
                        Login to Your Account
                      </a>
                    </td>
                  </tr>
                  <!-- Footer -->
                  <tr>
                    <td align="center" style="padding: 20px 0;">
                      <p style="font-size: 12px; color: #999;">&copy; 2025 NearbyAssist. All rights reserved.</p>
                    </td>
                  </tr>
                </table>
              </td>
            </tr>
          </table>
        </body>
        </html>
    `

	args := struct {
		Username  string
		Password  string
		LoginLink string
	}{
		Username:  p.Username,
		Password:  p.Password,
		LoginLink: p.LoginLink,
	}

	templ, err := template.New("invitation").Parse(emailTemplate)
	if err != nil {
		panic(err.Error())
	}

	buf := new(bytes.Buffer)
	if err := templ.Execute(buf, args); err != nil {
		panic(err.Error())
	}

	return buf.String()
}

func (p *PasswordResetPayload) GetRecipient() string {
	return p.Email
}

func (p *PasswordResetPayload) GetSubject() string {
	return "NearbyAssist Invitation Accepted"
}
