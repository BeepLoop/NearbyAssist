package config

import (
	"nearbyassist/internal/types"
	"os"

	"github.com/joho/godotenv"
)

var Env *types.Config

func Init() error {
	_, err := godotenv.Read()
	if err != nil {
		return err
	}

	Env = &types.Config{
		DSN: os.Getenv("DSN"),
	}

	return nil
}
