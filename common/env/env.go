package env

import (
	"os"

	"github.com/rs/zerolog"
)

func MustHaveEnv(logger zerolog.Logger, envName string) string {
	env, ok := os.LookupEnv(envName)
	if !ok {
		logger.Panic().Str("env", envName).Msg("Missing env")
	}
	return env
}
