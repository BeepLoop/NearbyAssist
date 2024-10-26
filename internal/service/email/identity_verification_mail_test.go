package email

import (
	"fmt"
	"nearbyassist/internal/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockMailService struct{}

func (m *mockMailService) SendMail(mail Mail) error {
	fmt.Println(mail)
	return nil
}

func TestIdentityVerificationMail(t *testing.T) {
	t.Run("Test setting To field", func(t *testing.T) {
		expected := []string{"example@email.com"}

		mailService := &mockMailService{}
		m := IdentityVerificationMail(mailService).To(expected)

		assert.Equal(t, expected, m.mail.To)
	})

	t.Run("Test sets correct content", func(t *testing.T) {
		mailService := &mockMailService{}

		m := IdentityVerificationMail(mailService)
		m.SetSubject(IDENTITY_REQUEST_ACKNOWLEDGMENT)

		assert.EqualValues(t, IDENTITY_REQUEST_ACKNOWLEDGMENT, m.mail.Subject)

		err := m.SetBody(IDENTITY_REQUEST_ACKNOWLEDGMENT, response.BasicEmailPayload{User: "foobar"})
		assert.NoError(t, err)

		m.To([]string{"jlmulit68@gmail.com"})
		assert.EqualValues(t, m.mail.To, []string{"jlmulit68@gmail.com"})

		err = m.Send()
		assert.NoError(t, err)
	})
}
