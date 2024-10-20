package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
