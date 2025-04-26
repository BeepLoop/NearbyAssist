package core

import (
	"slices"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	PASSWORD_MIN_LENGTH = 8
)

func BcryptPassword(pwd string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPwd), nil
}

func IsPasswordMatch(hashedPwd, plainPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)); err != nil {
		return false
	}

	return true
}

func IsPasswordSecure(password string) bool {
	allowed_special_chars := []string{"@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+", "-", "=", "!", "?"}

	if len(password) < PASSWORD_MIN_LENGTH {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecialChar := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else if slices.Contains(allowed_special_chars, string(char)) {
			hasSpecialChar = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecialChar
}
