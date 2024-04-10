package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN                      string
	Port                     string
	ApplicationProofLocation string
	ServicePhotoLocation     string
}

func LoadConfig() *Config {
	godotenv.Load()
	var dsn string

	if os.Getenv("GO_ENV") == "development" {
		dsn = os.Getenv("DSN_DEV")
	} else {
		dsn = os.Getenv("DSN_PROD")
	}

	return &Config{
		DSN:                      dsn,
		Port:                     os.Getenv("PORT"),
		ApplicationProofLocation: os.Getenv("APPLICATION_PROOF_LOCATION"),
		ServicePhotoLocation:     os.Getenv("SERVICE_PHOTO_LOCATION"),
	}
}
