package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVendorApplicationMail(t *testing.T) {
	t.Run("Test setting To field", func(t *testing.T) {
		expected := []string{"example@email.com"}

		mailService := &mockMailService{}
		m := VendorApplicationMail(mailService).To(expected)

		assert.Equal(t, expected, m.mail.To)
	})
}
