package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/VlasovArtem/hob/src/provider/service"
	"github.com/google/uuid"
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
	providerRouter := router.PathPrefix("/api/v1/providers").Subrouter()

	providerRouter.Path("").HandlerFunc(p.Add()).Methods("POST")
	providerRouter.Path("/{id}").HandlerFunc(p.Delete()).Methods("DELETE")
	providerRouter.Path("/{id}").HandlerFunc(p.Update()).Methods("PUT")
	providerRouter.Path("/{id}").HandlerFunc(p.FindById()).Methods("GET")
	providerRouter.Path("").
		Queries("userId", "{.*}").
		HandlerFunc(p.FindBy()).
		Methods("GET").
		Name("Find By")
}

func NewProviderHandler(providerService service.ProviderService) ProviderHandler {
	return &ProviderHandlerObject{providerService}
}

func (p *ProviderHandlerObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewProviderHandler(dependency.FindRequiredDependency[service.ProviderServiceObject, service.ProviderService](factory))
}

type ProviderHandler interface {
	Add() http.HandlerFunc
	Delete() http.HandlerFunc
	Update() http.HandlerFunc
	FindById() http.HandlerFunc
	FindBy() http.HandlerFunc
}

func (p *ProviderHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if body, err := rest.ReadRequestBody[model.CreateProviderRequest](request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Created(p.providerService.Add(body)).
				Perform()
		}
	}
}

func (p *ProviderHandlerObject) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)

			return
		} else {
			rest.NewAPIResponse(writer).
				NoContent(p.providerService.Delete(id)).
				Perform()
		}
	}
}

func (p *ProviderHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if body, err := rest.ReadRequestBody[model.UpdateProviderRequest](request); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				rest.NewAPIResponse(writer).
					Error(p.providerService.Update(id, body)).
					Perform()
			}
		}
	}
}

func (p *ProviderHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Ok(p.providerService.FindById(id)).
				Perform()
		}
	}
}

func (p *ProviderHandlerObject) FindBy() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId, err := rest.GetQueryParam[uuid.UUID](request, "userId")
		if err != nil {
			rest.HandleWithError(writer, err)
			return
		}

		page, err := rest.GetQueryParamOrDefault(request, "page", 0)
		if err != nil {
			rest.HandleWithError(writer, err)
			return
		}

		size, err := rest.GetQueryParamOrDefault(request, "size", 25)
		if err != nil {
			rest.HandleWithError(writer, err)
			return
		}

		name, err := rest.GetQueryParamOrDefault(request, "name", "")
		if err != nil {
			rest.HandleWithError(writer, err)
			return
		}

		rest.NewAPIResponse(writer).
			Ok(p.providerService.FindByNameLikeAndUserId(name, userId, page, size), nil).
			Perform()
	}
}
