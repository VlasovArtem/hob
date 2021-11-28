package environment

import (
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func GetEnvironmentVariable(name string, defaultValue string) string {
	if variable := os.Getenv(name); variable == "" {
		log.Info().Msgf("Environment variable with name '%s' not found, default used '%s'", name, defaultValue)
		return defaultValue
	} else {
		return variable
	}
}

func GetEnvironmentIntVariable(name string, defaultValue int) int {
	if variable := os.Getenv(name); variable == "" {
		log.Info().Msgf("Environment variable with name '%s' not found, default used '%d'", name, defaultValue)
		return defaultValue
	} else {
		if intVariable, err := strconv.Atoi(variable); err != nil {
			log.Fatal().Err(err)
			return 0
		} else {
			return intVariable
		}
	}
}
