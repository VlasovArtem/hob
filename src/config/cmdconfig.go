package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

const (
	API      = "a"
	TERMINAL = "t"
)

type CMDConfig struct {
	View string
}

func NewCMDConfig() *CMDConfig {
	return new(CMDConfig)
}

func (c *CMDConfig) ParseCMDConfig() {
	var value string
	pflag.StringVarP(&value, "view", "v", TERMINAL, "View: 't' - Terminal, 'a' - API. Default: terminal")
	pflag.Parse()

	if value != API && value != TERMINAL {
		log.Fatal().Msgf("View type %s is not supported. Possible values: 't' - Terminal, 'a' - API", value)
	}
	c.View = value
}

func (c *CMDConfig) IsTerminalView() bool {
	return c.View == TERMINAL
}
