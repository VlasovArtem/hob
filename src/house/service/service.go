package service

import (
	countries "country/service"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"house/model"
	"log"
)

type houseServiceObject struct {
	houses           map[uuid.UUID]model.House
	countriesService countries.CountryService
}

type HouseService interface {
	Add(house model.CreateHouseRequest) (model.HouseResponse, error)
	FindAll() []model.HouseResponse
	FindById(id uuid.UUID) (model.HouseResponse, error)
	ExistsById(id uuid.UUID) bool
}

func NewHouseService(countriesService countries.CountryService) HouseService {
	if countriesService == nil {
		log.Fatal("Country service is required")
	}

	return &houseServiceObject{
		houses:           make(map[uuid.UUID]model.House),
		countriesService: countriesService,
	}
}

func (h *houseServiceObject) Add(house model.CreateHouseRequest) (model.HouseResponse, error) {
	if err, country := h.countriesService.FindCountryByCode(house.Country); err != nil {
		return model.HouseResponse{}, err
	} else {
		house := house.ToEntity(&country)
		h.houses[house.Id] = house

		return house.ToResponse(), nil
	}
}

func (h *houseServiceObject) FindAll() []model.HouseResponse {
	result := make([]model.HouseResponse, 0)

	for _, house := range h.houses {
		result = append(result, house.ToResponse())
	}

	return result
}

func (h *houseServiceObject) FindById(id uuid.UUID) (model.HouseResponse, error) {
	if house, ok := h.houses[id]; ok {
		return house.ToResponse(), nil
	}

	return model.HouseResponse{}, errors.New(fmt.Sprintf("House with id - %s not exists", id))
}

func (h *houseServiceObject) ExistsById(id uuid.UUID) bool {
	_, ok := h.houses[id]

	return ok
}
