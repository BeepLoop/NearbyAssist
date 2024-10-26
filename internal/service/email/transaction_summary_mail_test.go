package email

import (
	"nearbyassist/internal/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionSummary(t *testing.T) {
	t.Run("Test setting To field", func(t *testing.T) {
		expected := []string{"example@email.com"}

		mailService := &mockMailService{}
		m := TransactionSummaryMail(mailService).To(expected)

		assert.Equal(t, expected, m.mail.To)
	})

	t.Run("Test setting html body", func(t *testing.T) {
		mailService := &mockMailService{}
		m := TransactionSummaryMail(mailService)

		data := response.TransactionSummary{}

		err := m.SetBody(data)
		assert.NoError(t, err)
	})
}
