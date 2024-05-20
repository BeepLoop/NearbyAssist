package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type AdminJwtClaims struct {
	AdminId  int       `json:"adminId"`
	Username string    `json:"username"`
	Role     AdminRole `json:"role"`
	jwt.RegisteredClaims
}

type UserJwtClaims struct {
	UserId int    `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
