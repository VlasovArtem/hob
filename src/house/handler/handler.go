package handler

import (
	"common/rest"
	"house/model"
	"house/service"
	"net/http"
)

type houseHandlerObject struct {
	houseService service.HouseService
}

func NewHouseHandler(houseService service.HouseService) HouseHandler {
	return &houseHandlerObject{houseService}
}

type HouseHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByUserId() http.HandlerFunc
}

func (h *houseHandlerObject) Add() http.HandlerFunc {
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

func (h *houseHandlerObject) FindById() http.HandlerFunc {
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

func (h *houseHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			responses := h.houseService.FindByUserId(id)

			rest.PerformResponse(writer, responses, nil)
		}
	}
}
