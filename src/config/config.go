package config

import (
	"fmt"
	"github.com/adrg/xdg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

var (
	DefaultConfigPath  = filepath.Join(getConfigDir(), "config.yml")
	DefaultLogFilePath = filepath.Join(os.TempDir(), fmt.Sprintf("hob.log"))
)

const (
	HOBConfig       = "HOBCONFIG"
	defaultView     = "terminal"
	DefaultLogLevel = "info"
)

type Config struct {
	User struct {
		Email    string
		Password string
	}
	App struct {
		LogFile  string
		LogLevel string
		View     string
	}
}

func NewConfig() *Config {
	return &Config{
		App: struct {
			LogFile  string
			LogLevel string
			View     string
		}{LogLevel: DefaultLogLevel, View: defaultView},
	}
}

func (c *Config) LoadConfig() {
	content, err := os.ReadFile(DefaultConfigPath)

	if err != nil {
		log.Warn().Err(err).Msg("Enable to load config file.")
		return
	}

	err = yaml.Unmarshal(content, c)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func (c *Config) EnrichWithCMD(cmdConfig *CMDConfig) {
	if cmdConfig.UserEmail != "" {
		c.User.Email = cmdConfig.UserEmail
	}
	if cmdConfig.UserPassword != "" {
		c.User.Password = cmdConfig.UserPassword
	}

	if cmdConfig.View != "" {
		c.App.View = cmdConfig.View
	} else if c.App.View != "" {
		c.App.View = defaultView
	}

	if cmdConfig.LogFile != "" {
		c.App.View = cmdConfig.LogFile
	} else if c.IsTerminalView() && c.App.LogFile == "" {
		c.App.LogFile = DefaultLogFilePath
	}

	if cmdConfig.LogLevel != "" {
		c.App.LogLevel = cmdConfig.LogLevel
	} else if c.App.LogLevel == "" {
		c.App.LogLevel = DefaultLogLevel
	}
}

func (c *Config) GetLogLevel() zerolog.Level {
	switch c.App.LogLevel {
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func getConfigDir() string {
	if env := os.Getenv(HOBConfig); env != "" {
		return env
	}

	if xdgHomePath, err := xdg.ConfigFile("hob"); err != nil {
		log.Fatal().Err(err)
		return ""
	} else {
		return xdgHomePath
	}
}

func (c *Config) IsTerminalView() bool {
	lowerView := strings.ToLower(c.App.View)
	return lowerView == defaultView || lowerView == "t"
}
