package utils

import (
	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/config"
	"nearbyassist/internal/models"
	"testing"
)

func TestGetUserIdFromJWT(t *testing.T) {
	tests := []struct {
		user           *models.UserModel
		expectedUserId int
	}{
		{
			user: &models.UserModel{
				Name:  "John Loyd Mulit",
				Email: "jlmulit68@gmail.com",
			},
			expectedUserId: 0,
		},
		{
			user: &models.UserModel{
				Model: models.Model{
					Id: 1,
				},
				Name:  "John Loyd Mulit",
				Email: "jlmulit68@gmail.com",
			},
			expectedUserId: 1,
		},
	}

	conf := config.LoadConfig()

	for _, test := range tests {
		jwtSigner := authenticator.NewJWTAuthenticator(conf)

		token, err := jwtSigner.GenerateUserAccessToken(test.user)
		if err != nil {
			t.Fatalf("Failed to create access token. error: %s", err.Error())
		}

		userId, err := GetUserIdFromJWT(jwtSigner, "Bearer "+token)
		if err != nil {
			t.Fatalf("Error: %s", err.Error())
		}

		if userId != test.expectedUserId {
			t.Fatalf("\nExpected: %v\nGot: %v\n", test.expectedUserId, userId)
		}
	}
}
