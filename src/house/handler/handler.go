package handler

import (
	"common/dependency"
	"common/rest"
	"github.com/gorilla/mux"
	"house/model"
	"house/service"
	"net/http"
)

type HouseHandlerObject struct {
	houseService service.HouseService
}

func NewHouseHandler(houseService service.HouseService) HouseHandler {
	return &HouseHandlerObject{houseService}
}

func (h *HouseHandlerObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewHouseHandler(factory.FindRequiredByObject(service.HouseServiceObject{}).(service.HouseService)),
	)
}

func (h *HouseHandlerObject) Init(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1/house").Subrouter()

	subrouter.Path("").HandlerFunc(h.Add()).Methods("POST")
	subrouter.Path("/{id}").HandlerFunc(h.FindById()).Methods("GET")
	subrouter.Path("/user/{id}").HandlerFunc(h.FindByUserId()).Methods("GET")
}

type HouseHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByUserId() http.HandlerFunc
}

func (h *HouseHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestEntity := model.CreateHouseRequest{}

		if err := rest.PerformRequest(&requestEntity, writer, request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			response, err := h.houseService.Add(requestEntity)

			rest.PerformResponseWithCode(writer, response, http.StatusCreated, err)
		}
	}
}

func (h *HouseHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			if house, err := h.houseService.FindById(id); err != nil {
				rest.HandleErrorResponseWithError(writer, http.StatusNotFound, err)
			} else {
				rest.PerformResponseWithCode(writer, house, http.StatusOK, nil)
			}
		}
	}
}

func (h *HouseHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			responses := h.houseService.FindByUserId(id)

			rest.PerformResponse(writer, responses, nil)
		}
	}
}
