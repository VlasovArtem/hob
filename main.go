package main

import (
	"github.com/VlasovArtem/hob/src/api"
	"github.com/VlasovArtem/hob/src/app"
	"github.com/VlasovArtem/hob/src/config"
	"github.com/VlasovArtem/hob/src/tui"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	cmdConfig := config.NewCMDConfig()
	cmdConfig.ParseCMDConfig()

	rootApplication := app.NewRootApplication()

	if cmdConfig.IsTerminalView() {
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
			Err(http.ListenAndServe(":3000", router)).
			Msg("HTTP Application error")
	}
}
