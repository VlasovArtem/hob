package service

import (
	countries "country/service"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"house/model"
	"log"
)

type house struct {
	houses           map[uuid.UUID]model.House
	countriesService countries.CountryService
}

type HouseService interface {
	AddHouse(house model.CreateHouseRequest) (error, model.HouseResponse)
	FindAllHouses() []model.HouseResponse
	FindById(id uuid.UUID) (error, model.HouseResponse)
}

func NewHouseService(countriesService countries.CountryService) HouseService {
	if countriesService == nil {
		log.Fatal("Country service is required")
	}

	return &house{
		houses:           make(map[uuid.UUID]model.House),
		countriesService: countriesService,
	}
}

func (h *house) AddHouse(house model.CreateHouseRequest) (error, model.HouseResponse) {
	if err, country := h.countriesService.FindCountryByCode(house.Country); err != nil {
		return err, model.HouseResponse{}
	} else {
		house := house.ToEntity(&country)
		h.houses[house.Id] = house

		return nil, house.ToResponse()
	}
}

func (h *house) FindAllHouses() []model.HouseResponse {
	result := []model.HouseResponse{}

	for _, house := range h.houses {
		result = append(result, house.ToResponse())
	}

	return result
}

func (h *house) FindById(id uuid.UUID) (error, model.HouseResponse) {
	if house, ok := h.houses[id]; ok {
		return nil, house.ToResponse()
	}

	return errors.New(fmt.Sprintf("House with id - %s not exists", id)), model.HouseResponse{}
}
