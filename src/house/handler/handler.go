package handler

import (
	"common/rest"
	helper "common/service"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
	AddHouseHandler() http.HandlerFunc
	FindAllHousesHandler() http.HandlerFunc
	FindHouseByIdHandler() http.HandlerFunc
}

func (h *houseHandlerObject) AddHouseHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestEntity := model.CreateHouseRequest{}

		if err := rest.PerformRequest(&requestEntity, writer, request); err == nil {
			if response, err := h.houseService.Add(requestEntity); err != nil {
				rest.HandleBadRequestWithError(writer, err)
			} else {
				writer.WriteHeader(http.StatusCreated)

				json.NewEncoder(writer).Encode(response)
			}
		} else {
			rest.HandleBadRequestWithError(writer, err)
		}
	}
}

func (h *houseHandlerObject) FindAllHousesHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		houses := h.houseService.FindAll()

		err := json.NewEncoder(writer).Encode(houses)

		if helper.HandleError(err, "Unable to encode response for find all request") {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (h *houseHandlerObject) FindHouseByIdHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		parameter, _ := rest.GetRequestParameter(request, "id")

		id, err := uuid.Parse(parameter)

		if err != nil {
			rest.HandleBadRequestWithError(writer, errors.New(fmt.Sprintf("the id is not valid %s", parameter)))

			return
		}

		if house, err := h.houseService.FindById(id); err != nil {
			rest.HandleErrorResponseWithError(writer, http.StatusNotFound, err)
		} else {
			content, err := json.Marshal(house)

			if err != nil {
				rest.HandleBadRequestWithError(writer, err)

				return
			}

			writer.Write(content)
		}

	}
}
