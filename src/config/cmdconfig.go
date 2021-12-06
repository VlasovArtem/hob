package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"strings"
)

var (
	viewParameters = map[string]*string{
		"terminal": nil,
		"t":        nil,
		"api":      nil,
		"a":        nil,
	}
)

type CMDConfig struct {
	UserEmail    string
	UserPassword string
	View         string
	LogFile      string
	LogLevel     string
}

func NewCMDConfig() *CMDConfig {
	return new(CMDConfig)
}

func (c *CMDConfig) ParseCMDConfig() {
	pflag.StringVarP(&c.View, "view", "v", "", "View: 't', 'terminal' - Terminal, 'a', 'api' - Api. Default: 'terminal', 't'")
	pflag.StringVarP(&c.LogFile, "log-file", "f", "", "Log file path. Default: Empty for api and temp for terminal.")
	pflag.StringVarP(&c.LogLevel, "log-level", "l", "", "Log level. Default: 'info'")
	pflag.StringVarP(&c.UserEmail, "user-email", "u", "", "Default user email")
	pflag.StringVarP(&c.UserEmail, "user-password", "p", "", "Default user password")
	pflag.Parse()

	c.validate()
}

func (c *CMDConfig) validate() {
	if c.View != "" {
		lowerView := strings.ToLower(c.View)
		if _, ok := viewParameters[lowerView]; !ok {
			log.Fatal().Msgf("View type %s is not supported. Possible values: 't' - Terminal, 'a' - Api", lowerView)
		}
	}
}
