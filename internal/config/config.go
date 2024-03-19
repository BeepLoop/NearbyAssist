package config

import (
	"nearbyassist/internal/types"
	"os"

	"github.com/joho/godotenv"
)

var Env *types.Config

func Init() error {
	godotenv.Load()
	var dsn string

	if os.Getenv("GO_ENV") == "development" {
		dsn = os.Getenv("DSN_DEV")
	} else {
		dsn = os.Getenv("DSN_PROD")
	}

	Env = &types.Config{
		DSN: dsn,
	}

	return nil
}
