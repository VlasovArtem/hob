package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	helper "helper/service"
	"house/model"
	"house/service"
	"net/http"
)

func AddHouseHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestEntity := model.CreateHouseRequest{}

		if err := helper.PerformRequest(&requestEntity, writer, request); err == nil {
			house := houseAddRequestToEntity(&requestEntity)

			service.AddHouse(house)

			writer.WriteHeader(http.StatusCreated)

			json.NewEncoder(writer).Encode(model.CreateHouseResponse{Id: house.Id})
		} else {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}
	}
}

func FindAllHousesHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		houses := service.FindAllHouses()

		err := json.NewEncoder(writer).Encode(houses)

		if helper.HandleError(err, "Unable to encode response for find all request") {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func FindHouseByIdHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		parameter, _ := helper.GetRequestParameter(request, "id")

		id, err := uuid.Parse(parameter)

		if err != nil {
			message := fmt.Sprintf("The id is not valid %s", parameter)
			http.Error(writer, message, http.StatusBadRequest)
			writer.Write([]byte(message))

			return
		}

		if err, house := service.FindById(id); err != nil {
			message := err.Error()
			http.Error(writer, message, http.StatusNotFound)
			writer.Write([]byte(message))
		} else {
			content, err := json.Marshal(house)

			if err != nil {
				message := err.Error()
				writer.Write([]byte(message))
				http.Error(writer, message, http.StatusBadRequest)

				return
			}

			writer.Write(content)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
		}

	}
}

func houseAddRequestToEntity(request *model.CreateHouseRequest) model.House {
	newUUID, _ := uuid.NewUUID()
	return model.House{
		Id:          newUUID,
		Name:        request.Name,
		Country:     request.Country,
		City:        request.City,
		StreetLine1: request.StreetLine1,
		StreetLine2: request.StreetLine2,
		Deleted:     true,
	}
}
