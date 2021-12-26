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

func (h *HouseHandlerObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewHouseHandler(factory.FindRequiredByObject(service.HouseServiceObject{}).(service.HouseService))
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
		requestEntity := model.CreateHouseRequest{}

		if err := rest.ReadRequestBody(&requestEntity, writer, request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			response, err := h.houseService.Add(requestEntity)

			rest.PerformResponseWithCode(writer, response, http.StatusCreated, err)
		}
	}
}

func (h *HouseHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if house, err := h.houseService.FindById(id); err != nil {
				rest.HandleErrorResponseWithError(writer, http.StatusNotFound, err)
			} else {
				rest.PerformResponseWithCode(writer, house, http.StatusOK, nil)
			}
		}
	}
}

func (h *HouseHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			updateHouseRequest := model.UpdateHouseRequest{}

			if err := rest.ReadRequestBody(&updateHouseRequest, writer, request); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				rest.PerformResponseWithBody(writer, nil, h.houseService.Update(id, updateHouseRequest))
			}
		}
	}
}

func (h *HouseHandlerObject) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.PerformResponseWithCode(writer, nil, http.StatusNoContent, h.houseService.DeleteById(id))
		}
	}
}

func (h *HouseHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			responses := h.houseService.FindByUserId(id)

			rest.PerformResponseWithBody(writer, responses, nil)
		}
	}
}
