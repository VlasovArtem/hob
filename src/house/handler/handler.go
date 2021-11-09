package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	helper "helper/service"
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

		if err := helper.PerformRequest(&requestEntity, writer, request); err == nil {
			if err, response := h.houseService.AddHouse(requestEntity); err != nil {
				helper.HandleBadRequestWithError(writer, err)
			} else {
				writer.WriteHeader(http.StatusCreated)

				json.NewEncoder(writer).Encode(response)
			}
		} else {
			helper.HandleBadRequestWithError(writer, err)
		}
	}
}

func (h *houseHandlerObject) FindAllHousesHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		houses := h.houseService.FindAllHouses()

		err := json.NewEncoder(writer).Encode(houses)

		if helper.HandleError(err, "Unable to encode response for find all request") {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (h *houseHandlerObject) FindHouseByIdHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		parameter, _ := helper.GetRequestParameter(request, "id")

		id, err := uuid.Parse(parameter)

		if err != nil {
			helper.HandleBadRequestWithError(writer, errors.New(fmt.Sprintf("the id is not valid %s", parameter)))

			return
		}

		if err, house := h.houseService.FindById(id); err != nil {
			helper.HandleErrorResponseWithError(writer, http.StatusNotFound, err)
		} else {
			content, err := json.Marshal(house)

			if err != nil {
				helper.HandleBadRequestWithError(writer, err)

				return
			}

			writer.Write(content)
		}

	}
}
