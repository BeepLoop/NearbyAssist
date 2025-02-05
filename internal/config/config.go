package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	Instance *Config
)

type Config struct {
	PORT            string
	LOG_FILE        string
	ALLOWED_ORIGINS []string

	DB_USER string
	DB_PWD  string
	DB_NAME string
	DB_HOST string
	DB_PORT string
	DB_NET  string

	JWT_SECRET   string
	JWT_DURATION int

	ENCRYPTION_KEY string

	ROUTE_ENGINE_URL string

	ONE_SIGNAL_APP_ID  string
	ONE_SIGNAL_API_KEY string

	APPLICATION_PROOF_DIR string
	POLICE_CLEARANCE_DIR  string
	SERVICE_PHOTO_DIR     string
	SYS_COMPLAINT_DIR     string
	ID_FRONT_DIR          string
	ID_BACK_DIR           string
	FACE_IMG_DIR          string
}

func GetConfig() *Config {
	if Instance != nil {
		return Instance
	}

	Instance = initialize()
	return Instance
}

func initialize() *Config {
	godotenv.Load()

	jwtDuration := getEnv("JWT_DURATION", "60")
	duration, err := strconv.Atoi(jwtDuration)
	if err != nil {
		panic("JWT_DURATION must be an integer value (in seconds)")
	}

	return &Config{
		PORT:            getEnv("PORT", "3000"),
		LOG_FILE:        getEnv("LOG_FILE", "logs/server.log"),
		ALLOWED_ORIGINS: strings.Split(getEnv("ALLOWED_ORIGINS", "http://127.0.0.1:3001"), ","),

		DB_USER: getEnv("DB_USER", "root"),
		DB_PWD:  getEnv("DB_PASSWORD", "secret"),
		DB_NAME: getEnv("DB_NAME", "nearbyassist"),
		DB_HOST: getEnv("DB_HOST", "127.0.0.1"),
		DB_PORT: getEnv("DB_PORT", "3306"),
		DB_NET:  getEnv("DB_NET", "tcp"),

		JWT_SECRET:   getEnv("JWT_SECRET", "secret"),
		JWT_DURATION: duration,

		ENCRYPTION_KEY: getEnv("ENCRYPTION_KEY", "key"),

		APPLICATION_PROOF_DIR: getEnv("APPLICATION_PROOF_DIR", "uploads/application_proof"),
		POLICE_CLEARANCE_DIR:  getEnv("POLICE_CLEARANCE_DIR", "uploads/police_clearance"),
		SERVICE_PHOTO_DIR:     getEnv("SERVICE_PHOTO_DIR", "uploads/service_photo"),
		SYS_COMPLAINT_DIR:     getEnv("SYSTEM_COMPLAINT_DIR", "uploads/system_complaint"),
		ID_FRONT_DIR:          getEnv("VERIFICATION_FRONT_ID_DIR", "uploads/verification/front_id"),
		ID_BACK_DIR:           getEnv("VERIFICATION_BACK_ID_DIR", "uploads/verification/back_id"),
		FACE_IMG_DIR:          getEnv("VERIFICATION_FACE_DIR", "uploads/verification/face"),

		ROUTE_ENGINE_URL: getEnv("ROUTE_ENGINE_URL", "http://127.0.0.1:5000"),

		ONE_SIGNAL_APP_ID:  mustGetEnv("ONE_SIGNAL_APP_ID"),
		ONE_SIGNAL_API_KEY: mustGetEnv("ONE_SIGNAL_API_KEY"),
	}
}

func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	return value
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("Environment key '%s' does not exits in environment", key))
	}

	return value
}
