package email

import (
	"nearbyassist/internal/response"
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

	t.Run("Test setting html body", func(t *testing.T) {
		mailService := &mockMailService{}

		m := VendorApplicationMail(mailService)

		data := response.BasicEmailPayload{
			User:            "Jane Doe",
			SupportEndpoint: "http://localhost:3000/support",
		}

		m.To([]string{"jlmulit68@gmail.com"})
		assert.EqualValues(t, []string{"jlmulit68@gmail.com"}, m.mail.To)

		m.SetSubject(VENDOR_APPLICATION_ACKNOWLEDGMENT)

		err := m.SetBody(VENDOR_APPLICATION_ACKNOWLEDGMENT, data)
		assert.NoError(t, err)
	})
}
