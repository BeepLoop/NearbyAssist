package config

import (
	"github.com/joho/godotenv"
	"nearbyassist/internal/types"
)

var Env *types.Config

func Init() error {
	envMap, err := godotenv.Read()
	if err != nil {
		return err
	}

	Env = &types.Config{
		DSN: envMap["DSN"],
	}

	return nil
}
