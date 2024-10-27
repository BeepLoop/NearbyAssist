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

	t.Run("Test setting html body", func(t *testing.T) {
		mailService := &mockMailService{}

		m := IdentityVerificationMail(mailService)

		data := response.BasicEmailPayload{
			User:            "Jane Doe",
			SupportEndpoint: "http://localhost:3000/support",
		}

		m.To([]string{"jlmulit68@gmail.com"})
		assert.EqualValues(t, []string{"jlmulit68@gmail.com"}, m.mail.To)

		m.SetSubject(IDENTITY_REQUEST_ACKNOWLEDGMENT)

		err := m.SetBody(IDENTITY_REQUEST_ACKNOWLEDGMENT, data)
		assert.NoError(t, err)
	})
}
