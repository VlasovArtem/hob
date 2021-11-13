package service

import (
	countries "country/service"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"house/model"
	"log"
	"user/service"
)

type houseServiceObject struct {
	houses           map[uuid.UUID]model.House
	userHouses       map[uuid.UUID][]model.House
	countriesService countries.CountryService
	userService      service.UserService
}

type HouseService interface {
	Add(house model.CreateHouseRequest) (model.HouseResponse, error)
	FindById(id uuid.UUID) (model.HouseResponse, error)
	FindByUserId(userId uuid.UUID) []model.HouseResponse
	ExistsById(id uuid.UUID) bool
}

func NewHouseService(
	countriesService countries.CountryService,
	userService service.UserService,
) HouseService {
	if countriesService == nil {
		log.Fatal("Country service is required")
	}

	return &houseServiceObject{
		houses:           make(map[uuid.UUID]model.House),
		userHouses:       make(map[uuid.UUID][]model.House),
		countriesService: countriesService,
		userService:      userService,
	}
}

func (h *houseServiceObject) Add(request model.CreateHouseRequest) (response model.HouseResponse, err error) {
	if err, country := h.countriesService.FindCountryByCode(request.Country); err != nil {
		return response, err
	} else if !h.userService.ExistsById(request.UserId) {
		return response, errors.New(fmt.Sprintf("user with id %s in not exists", request.UserId))
	} else {
		entity := request.ToEntity(&country)
		h.houses[entity.Id] = entity
		h.userHouses[entity.UserId] = append(h.userHouses[entity.UserId], entity)

		return entity.ToResponse(), nil
	}
}

func (h *houseServiceObject) FindById(id uuid.UUID) (model.HouseResponse, error) {
	if house, ok := h.houses[id]; ok {
		return house.ToResponse(), nil
	}

	return model.HouseResponse{}, errors.New(fmt.Sprintf("house with id %s not found", id))
}

func (h *houseServiceObject) FindByUserId(userId uuid.UUID) []model.HouseResponse {
	if houses, ok := h.userHouses[userId]; ok {
		var response []model.HouseResponse

		for _, house := range houses {
			response = append(response, house.ToResponse())
		}

		return response
	}
	return []model.HouseResponse{}
}

func (h *houseServiceObject) ExistsById(id uuid.UUID) bool {
	_, ok := h.houses[id]

	return ok
}
