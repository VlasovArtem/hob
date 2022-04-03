package handler

import (
	"encoding/json"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

type HealthHandlerObject struct{}

func NewHealthHandler() *HealthHandlerObject {
	return &HealthHandlerObject{}
}

func (h *HealthHandlerObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewHealthHandler()
}

func (h *HealthHandlerObject) Init(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1/health").Subrouter()

	subrouter.Path("").HandlerFunc(h.HealthCheck()).Methods("GET")
}

type HealthHandler interface {
	HealthCheck() http.HandlerFunc
}

func (h *HealthHandlerObject) HealthCheck() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		marshal, err := json.Marshal(HealthStatus{
			status: "UP",
		})

		if err != nil {
			log.Error().Err(err)
		}

		_, err = writer.Write(marshal)

		if err != nil {
			log.Error().Err(err)
		}
	}
}

type HealthStatus struct {
	status string
}
