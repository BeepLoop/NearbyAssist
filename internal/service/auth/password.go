package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func IsPasswordMatch(hashedPwd, plainPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)); err != nil {
		return false
	}
	return true
}
