package email

import (
	"nearbyassist/internal/response"
	"testing"
	"time"

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

		data := response.TransactionSummary{
			Id:                   "Sample transaction ID",
			Vendor:               "Jane Doe",
			Client:               "John Doe",
			Price:                "1000",
			CreatedAt:            time.Now().Format(time.RFC822Z),
			StartDate:            time.Now().Format(time.RFC822Z),
			EndDate:              time.Now().Format(time.RFC822Z),
			ServiceTitle:         "Sample service title",
			VendorEmail:          "provider@email.com",
			ClientEmail:          "client@email.com",
			ConfirmationEndpoint: "http://localhost:3000/confirm",
		}

		m.To([]string{"jlmulit68@gmail.com"})
		assert.EqualValues(t, []string{"jlmulit68@gmail.com"}, m.mail.To)

		err := m.SetBody(data)
		assert.NoError(t, err)
	})
}
