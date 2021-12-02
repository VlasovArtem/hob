package config

import (
	"github.com/adrg/xdg"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

type XDGConfig struct {
	User struct {
		Email    string
		Password string
	}
}

func ReadXDGConfig() *XDGConfig {
	configFile, err := xdg.ConfigFile("hob/config.yml")

	if err != nil {
		log.Warn().Err(err)
		return nil
	}

	content, err := os.ReadFile(configFile)

	if err != nil {
		log.Warn().Err(err)
		return nil
	}

	config := XDGConfig{}

	err = yaml.Unmarshal(content, &config)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	return &config
}
