package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/house/service"
	"github.com/gorilla/mux"
	"net/http"
)

type HouseHandlerObject struct {
	houseService service.HouseService
}

func NewHouseHandler(houseService service.HouseService) HouseHandler {
	return &HouseHandlerObject{houseService}
}

func (h *HouseHandlerObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewHouseHandler(dependency.FindRequiredDependency[service.HouseServiceObject, service.HouseService](factory))
}

func (h *HouseHandlerObject) Init(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1/house").Subrouter()

	subrouter.Path("").HandlerFunc(h.Add()).Methods("POST")
	subrouter.Path("/{id}").HandlerFunc(h.FindById()).Methods("GET")
	subrouter.Path("/{id}").HandlerFunc(h.Update()).Methods("DELETE")
	subrouter.Path("/{id}").HandlerFunc(h.Delete()).Methods("PUT")
	subrouter.Path("/user/{id}").HandlerFunc(h.FindByUserId()).Methods("GET")
}

type HouseHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
	FindByUserId() http.HandlerFunc
}

func (h *HouseHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestBody, err := rest.ReadRequestBody[model.CreateHouseRequest](request)
		if err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Created(h.houseService.Add(requestBody)).
				Perform()
		}
	}
}

func (h *HouseHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Ok(h.houseService.FindById(id)).
				Perform()
		}
	}
}

func (h *HouseHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if body, err := rest.ReadRequestBody[model.UpdateHouseRequest](request); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				rest.NewAPIResponse(writer).
					Error(h.houseService.Update(id, body)).
					Perform()
			}
		}
	}
}

func (h *HouseHandlerObject) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				StatusCode(http.StatusNoContent).
				Error(h.houseService.DeleteById(id)).
				Perform()
		}
	}
}

func (h *HouseHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Body(h.houseService.FindByUserId(id)).
				Perform()
		}
	}
}
