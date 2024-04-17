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
	DB_User                  string
	DB_Password              string
	DB_Name                  string
	DB_Host                  string
	DB_Port                  string
	Port                     string
	JwtSecret                string
	StorageType              StorageType
	ApplicationProofLocation string
	ServicePhotoLocation     string
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		DB_User:                  os.Getenv("DB_USER"),
		DB_Password:              os.Getenv("DB_PASSWORD"),
		DB_Name:                  os.Getenv("DB_NAME"),
		DB_Host:                  os.Getenv("DB_HOST"),
		DB_Port:                  os.Getenv("DB_PORT"),
		Port:                     os.Getenv("PORT"),
		JwtSecret:                os.Getenv("JWT_SECRET"),
		StorageType:              StorageType(os.Getenv("STORAGE_TYPE")),
		ApplicationProofLocation: os.Getenv("APPLICATION_PROOF_LOCATION"),
		ServicePhotoLocation:     os.Getenv("SERVICE_PHOTO_LOCATION"),
	}
}
