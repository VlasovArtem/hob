package main

import (
	"github.com/VlasovArtem/hob/src/api"
	"github.com/VlasovArtem/hob/src/app"
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/config"
	"github.com/VlasovArtem/hob/src/tui"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

func main() {
	cfg := prepareConfig()

	if cfg.App.LogFile != "" {
		common.EnsurePath(cfg.App.LogFile, common.DefaultDirMod)

		mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
		file, err := os.OpenFile(cfg.App.LogFile, mod, common.DefaultFileMode)

		defer func() {
			_ = file.Close()
		}()

		if err != nil {
			log.Fatal().Err(err)
		}

		log.Logger = log.Output(zerolog.ConsoleWriter{Out: file})
	}

	zerolog.SetGlobalLevel(cfg.GetLogLevel())

	rootApplication := app.NewRootApplication(cfg)

	startApplication(cfg, rootApplication)
}

func startApplication(cfg *config.Config, rootApplication *app.RootApplication) {
	if cfg.IsTerminalView() {
		tApp := tui.NewTApp(rootApplication)

		tApp.Init()

		if err := tApp.Run(); err != nil {
			panic(err)
		}
	} else {
		router := mux.NewRouter().StrictSlash(true)

		api.InitApi(router, rootApplication)

		http.Handle("/", router)

		if err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			if template, err := route.GetPathTemplate(); err != nil {
				log.Error().Err(err)
			} else {
				log.Info().Msg(template)
			}
			return nil
		}); err != nil {
			log.Fatal().Err(err).Msg("router walk error")
		}

		log.Fatal().
			Err(http.ListenAndServe(":3030", router)).
			Msg("HTTP Application error")
	}
}

func prepareConfig() *config.Config {
	cmdConfig := config.NewCMDConfig()
	cmdConfig.ParseCMDConfig()

	cfg := config.NewConfig()
	cfg.LoadConfig()
	cfg.EnrichWithCMD(cmdConfig)

	return cfg
}
