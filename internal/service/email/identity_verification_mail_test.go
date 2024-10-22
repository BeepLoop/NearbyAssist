package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockMailService struct{}

func (m *mockMailService) SendMail(mail Mail) error {
	return nil
}

func TestIdentityVerificationMail(t *testing.T) {
	t.Run("Test setting To field", func(t *testing.T) {
		expected := []string{"example@email.com"}

		mailService := &mockMailService{}
		m := IdentityVerificationMail(mailService).To(expected)

		assert.Equal(t, expected, m.mail.To)
	})
}
