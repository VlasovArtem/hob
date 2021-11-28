package handler

import (
	"encoding/json"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/VlasovArtem/hob/src/provider/service"
	"github.com/gorilla/mux"
	"net/http"
)

type ProviderHandlerObject struct {
	providerService service.ProviderService
}

func (p *ProviderHandlerObject) Init(router *mux.Router) {
	providerRouter := router.PathPrefix("/api/v1/provider").Subrouter()

	providerRouter.Path("").HandlerFunc(p.Add()).Methods("POST")
	providerRouter.Path("/{id}").HandlerFunc(p.FindById()).Methods("GET")
	providerRouter.Path("/name/{name}").HandlerFunc(p.FindByNameLike()).Methods("GET")
}

func NewProviderHandler(providerService service.ProviderService) ProviderHandler {
	return &ProviderHandlerObject{providerService}
}

func (p *ProviderHandlerObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(NewProviderHandler(factory.FindRequiredByObject(service.ProviderServiceObject{}).(service.ProviderService)))
}

type ProviderHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByNameLike() http.HandlerFunc
}

func (p *ProviderHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestEntity := model.CreateProviderRequest{}

		if err := rest.PerformRequest(&requestEntity, writer, request); err != nil {
			rest.HandleBadRequestWithError(writer, err)

			return
		}
		dto, err := p.providerService.Add(requestEntity)
		rest.PerformResponseWithCode(writer, dto, http.StatusCreated, err)
	}
}

func (p *ProviderHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)

			return
		} else {
			if dto, err := p.providerService.FindById(id); err != nil {
				rest.HandleErrorResponseWithError(writer, http.StatusNotFound, err)
			} else {
				json.NewEncoder(writer).Encode(dto)
			}
		}
	}
}

func (p *ProviderHandlerObject) FindByNameLike() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		name, err := rest.GetRequestParameter(request, "name")
		if err != nil {
			rest.HandleBadRequestWithError(writer, err)
			return
		}

		page, err := rest.GetQueryIntParameterOrDefault(request, "page", 0)
		if err != nil {
			rest.HandleBadRequestWithError(writer, err)
			return
		}

		size, err := rest.GetQueryIntParameterOrDefault(request, "size", 25)
		if err != nil {
			rest.HandleBadRequestWithError(writer, err)
			return
		}

		rest.PerformResponse(writer, p.providerService.FindByNameLike(name, page, size), nil)
	}
}
