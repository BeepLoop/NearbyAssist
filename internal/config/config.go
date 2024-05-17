package config

import (
	"os"
	"strconv"

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
	JwtDuration              int
	StorageType              StorageType
	ApplicationProofLocation string
	ServicePhotoLocation     string
	RouteEngineUrl           string
}

func LoadConfig() *Config {
	godotenv.Load()

	jwtDuration := os.Getenv("JWT_DURATION")
	duration, err := strconv.Atoi(jwtDuration)
	if err != nil {
		panic("JWT_DURATION must be an integer value (in seconds)")
	}

	return &Config{
		DB_User:                  os.Getenv("DB_USER"),
		DB_Password:              os.Getenv("DB_PASSWORD"),
		DB_Name:                  os.Getenv("DB_NAME"),
		DB_Host:                  os.Getenv("DB_HOST"),
		DB_Port:                  os.Getenv("DB_PORT"),
		Port:                     os.Getenv("PORT"),
		JwtSecret:                os.Getenv("JWT_SECRET"),
		JwtDuration:              duration,
		StorageType:              StorageType(os.Getenv("STORAGE_TYPE")),
		ApplicationProofLocation: os.Getenv("APPLICATION_PROOF_LOCATION"),
		ServicePhotoLocation:     os.Getenv("SERVICE_PHOTO_LOCATION"),
		RouteEngineUrl:           os.Getenv("ROUTE_ENGINE_URL"),
	}
}
