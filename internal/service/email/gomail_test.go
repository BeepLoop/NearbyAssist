package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmail(t *testing.T) {
	t.Run("Test empty params on constructor", func(t *testing.T) {
		email := ""
		password := ""

		_, err := NewGoMail(email, password)
		assert.Error(t, err)
	})
}
