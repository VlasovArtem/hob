package handler

import (
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

type FindByNameRequest struct {
	Name string
}

func (p *ProviderHandlerObject) Init(router *mux.Router) {
	providerRouter := router.PathPrefix("/api/v1/provider").Subrouter()

	providerRouter.Path("").HandlerFunc(p.Add()).Methods("POST")
	providerRouter.Path("/{id}").HandlerFunc(p.Delete()).Methods("DELETE")
	providerRouter.Path("/{id}").HandlerFunc(p.Update()).Methods("PUT")
	providerRouter.Path("/{id}").HandlerFunc(p.FindById()).Methods("GET")
	providerRouter.Path("/user/{id}").HandlerFunc(p.FindByUserId()).Methods("GET")
	providerRouter.Path("/user/{id}").HandlerFunc(p.FindByNameLikeAndUserId()).Methods("POST")
}

func NewProviderHandler(providerService service.ProviderService) ProviderHandler {
	return &ProviderHandlerObject{providerService}
}

func (p *ProviderHandlerObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewProviderHandler(factory.FindRequiredByObject(service.ProviderServiceObject{}).(service.ProviderService))
}

type ProviderHandler interface {
	Add() http.HandlerFunc
	Delete() http.HandlerFunc
	Update() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByNameLikeAndUserId() http.HandlerFunc
	FindByUserId() http.HandlerFunc
}

func (p *ProviderHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestEntity := model.CreateProviderRequest{}

		if err := rest.ReadRequestBody(&requestEntity, writer, request); err != nil {
			rest.HandleWithError(writer, err)

			return
		}
		dto, err := p.providerService.Add(requestEntity)
		rest.PerformResponseWithCode(writer, dto, http.StatusCreated, err)
	}
}

func (p *ProviderHandlerObject) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)

			return
		} else {
			dto, err := p.providerService.FindById(id)

			rest.PerformResponseWithBody(writer, dto, err)
		}
	}
}

func (p *ProviderHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)

			return
		} else {
			requestEntity := model.UpdateProviderRequest{}

			if err = rest.ReadRequestBody(&requestEntity, writer, request); err != nil {
				rest.HandleWithError(writer, err)

				return
			}
			rest.PerformResponseWithBody(writer, nil, p.providerService.Update(id, requestEntity))
		}
	}
}

func (p *ProviderHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)

			return
		} else {
			dto, err := p.providerService.FindById(id)

			rest.PerformResponseWithBody(writer, dto, err)
		}
	}
}

func (p *ProviderHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := rest.GetIdRequestParameter(request)

		if err != nil {
			rest.HandleWithError(writer, err)

			return
		}

		rest.PerformResponseWithBody(writer, p.providerService.FindByUserId(id), nil)
	}
}

func (p *ProviderHandlerObject) FindByNameLikeAndUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestEntity := FindByNameRequest{}

		if err := rest.ReadRequestBody(&requestEntity, writer, request); err != nil {

			return
		}

		page, err := rest.GetQueryIntParameterOrDefault(request, "page", 0)
		if err != nil {
			rest.HandleWithError(writer, err)
			return
		}

		size, err := rest.GetQueryIntParameterOrDefault(request, "size", 25)
		if err != nil {
			rest.HandleWithError(writer, err)
			return
		}

		id, err := rest.GetIdRequestParameter(request)

		if err != nil {
			rest.HandleWithError(writer, err)

			return
		}

		rest.PerformResponseWithBody(writer, p.providerService.FindByNameLikeAndUserId(requestEntity.Name, id, page, size), nil)
	}
}
