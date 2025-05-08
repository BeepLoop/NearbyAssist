package config

import (
	"fmt"
	"nearbyassist/internal/utils"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	Instance *Config
)

type Config struct {
	DOMAIN string

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
	ID_FRONT_DIR          string
	ID_BACK_DIR           string
	FACE_IMG_DIR          string
	BUG_REPORT_DIR        string
	REPORT_USER_DIR       string

	PDF_GENERATION_TIMEOUT int
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

	jwtDuration := GetEnv("JWT_DURATION", "60")
	duration, err := strconv.Atoi(jwtDuration)
	if err != nil {
		panic("JWT_DURATION must be an integer value (in seconds)")
	}

	return &Config{
		DOMAIN: GetEnv("DOMAIN", "http://localhost:3000"),

		PORT:            GetEnv("PORT", "3000"),
		LOG_FILE:        GetEnv("LOG_FILE", "logs/server.log"),
		ALLOWED_ORIGINS: strings.Split(GetEnv("ALLOWED_ORIGINS", "http://127.0.0.1:3001"), ","),

		DB_USER: GetEnv("DB_USER", "root"),
		DB_PWD:  GetEnv("DB_PASSWORD", "secret"),
		DB_NAME: GetEnv("DB_NAME", "nearbyassist"),
		DB_HOST: GetEnv("DB_HOST", "127.0.0.1"),
		DB_PORT: GetEnv("DB_PORT", "3306"),
		DB_NET:  GetEnv("DB_NET", "tcp"),

		JWT_SECRET:   GetEnv("JWT_SECRET", "secret"),
		JWT_DURATION: duration,

		ENCRYPTION_KEY: GetEnv("ENCRYPTION_KEY", "key"),

		APPLICATION_PROOF_DIR: GetEnv("APPLICATION_PROOF_DIR", "uploads/application_proof"),
		POLICE_CLEARANCE_DIR:  GetEnv("POLICE_CLEARANCE_DIR", "uploads/police_clearance"),
		SERVICE_PHOTO_DIR:     GetEnv("SERVICE_PHOTO_DIR", "uploads/service_photo"),
		ID_FRONT_DIR:          GetEnv("VERIFICATION_FRONT_ID_DIR", "uploads/verification/front_id"),
		ID_BACK_DIR:           GetEnv("VERIFICATION_BACK_ID_DIR", "uploads/verification/back_id"),
		FACE_IMG_DIR:          GetEnv("VERIFICATION_FACE_DIR", "uploads/verification/face"),
		BUG_REPORT_DIR:        GetEnv("BUG_REPORT_DIR", "uploads/bug_report"),
		REPORT_USER_DIR:       GetEnv("REPORT_USER_DIR", "uploads/report_user"),

		ROUTE_ENGINE_URL: GetEnv("ROUTE_ENGINE_URL", "http://127.0.0.1:5000"),

		ONE_SIGNAL_APP_ID:  MustGetEnv("ONE_SIGNAL_APP_ID"),
		ONE_SIGNAL_API_KEY: MustGetEnv("ONE_SIGNAL_API_KEY"),

		PDF_GENERATION_TIMEOUT: utils.Must(strconv.Atoi(GetEnv("PDF_GENERATION_TIMEOUT", "120"))),
	}
}

func GetEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	return value
}

func MustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("Environment key '%s' does not exits in environment", key))
	}

	return value
}
