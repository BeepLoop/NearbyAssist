package config

import (
	"nearbyassist/internal/types"
	"os"

	"github.com/labstack/gommon/log"
)

var Env *types.Config

func Init() error {
	Env = &types.Config{
		DSN: os.Getenv("DSN"),
	}

	dsn := os.Getenv("DSN")
	log.Info("DSN: ", dsn)

	return nil
}
