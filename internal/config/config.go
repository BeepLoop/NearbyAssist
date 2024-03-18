package config

import (
	"nearbyassist/internal/types"
	"os"
)

var Env *types.Config

func Init() error {
	Env = &types.Config{
		DSN: os.Getenv("DSN"),
	}

	return nil
}
