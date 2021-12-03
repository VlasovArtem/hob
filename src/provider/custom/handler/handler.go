package handler

import (
	"encoding/json"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/provider/custom/model"
	"github.com/VlasovArtem/hob/src/provider/custom/service"
	"github.com/gorilla/mux"
	"net/http"
)

type CustomProviderHandlerObject struct {
	customProviderService service.CustomProviderService
}

func (c *CustomProviderHandlerObject) Init(router *mux.Router) {
	customProviderRouter := router.PathPrefix("/api/v1/provider/custom").Subrouter()

	customProviderRouter.Path("").HandlerFunc(c.Add()).Methods("POST")
	customProviderRouter.Path("/{id}").HandlerFunc(c.FindById()).Methods("GET")
	customProviderRouter.Path("/user/{id}").HandlerFunc(c.FindByUserId()).Methods("GET")
}

func NewCustomProviderHandler(customProviderService service.CustomProviderService) CustomProviderHandler {
	return &CustomProviderHandlerObject{customProviderService}
}

func (c *CustomProviderHandlerObject) Initialize(factory dependency.DependenciesFactory) interface{} {
	return NewCustomProviderHandler(factory.FindRequiredByObject(service.CustomProviderServiceObject{}).(service.CustomProviderService))
}

type CustomProviderHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByUserId() http.HandlerFunc
}

func (c *CustomProviderHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestEntity := model.CreateCustomProviderRequest{}

		if err := rest.PerformRequest(&requestEntity, writer, request); err != nil {
			rest.HandleBadRequestWithError(writer, err)

			return
		}
		dto, err := c.customProviderService.Add(requestEntity)
		rest.PerformResponseWithCode(writer, dto, http.StatusCreated, err)
	}
}

func (c *CustomProviderHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)

			return
		} else {
			if dto, err := c.customProviderService.FindById(id); err != nil {
				rest.HandleErrorResponseWithError(writer, http.StatusNotFound, err)
			} else {
				json.NewEncoder(writer).Encode(dto)
			}
		}
	}
}

func (c *CustomProviderHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)

			return
		} else {
			json.NewEncoder(writer).Encode(c.customProviderService.FindByUserId(id))
		}
	}
}
