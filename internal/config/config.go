package config

import (
	"os"

	"github.com/joho/godotenv"
)

type StorageType string

const (
	STORAGE_DISK  StorageType = "disk"
	STORAGE_DUMMY StorageType = "dummy"
)

type Config struct {
	DSN                      string
	Port                     string
	JwtSecret                string
	StorageType              StorageType
	ApplicationProofLocation string
	ServicePhotoLocation     string
	RouteEngineUrl           string
}

func LoadConfig() *Config {
	godotenv.Load()

	environment := os.Getenv("GO_ENV")

	var dsn string
	if environment == "development" {
		dsn = os.Getenv("DSN_DEV")
	} else {
		dsn = os.Getenv("DSN_PROD")
	}

	return &Config{
		DSN:                      dsn,
		Port:                     os.Getenv("PORT"),
		JwtSecret:                os.Getenv("JWT_SECRET"),
		StorageType:              StorageType(os.Getenv("STORAGE_TYPE")),
		ApplicationProofLocation: os.Getenv("APPLICATION_PROOF_LOCATION"),
		ServicePhotoLocation:     os.Getenv("SERVICE_PHOTO_LOCATION"),
		RouteEngineUrl:           os.Getenv("ROUTE_ENGINE_URL"),
	}
}
