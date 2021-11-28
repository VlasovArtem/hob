package service

import (
	"common/dependency"
	countries "country/service"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"house/model"
	"house/respository"
	"log"
	userService "user/service"
)

type HouseServiceObject struct {
	countriesService countries.CountryService
	userService      userService.UserService
	repository       respository.HouseRepository
}

func NewHouseService(
	countriesService countries.CountryService,
	userService userService.UserService,
	repository respository.HouseRepository,
) HouseService {
	if countriesService == nil {
		log.Fatal("Country service is required")
	}

	return &HouseServiceObject{
		countriesService: countriesService,
		userService:      userService,
		repository:       repository,
	}
}

func (h *HouseServiceObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewHouseService(
			factory.FindRequiredByObject(countries.CountryServiceObject{}).(countries.CountryService),
			factory.FindRequiredByObject(userService.UserServiceObject{}).(userService.UserService),
			factory.FindRequiredByObject(respository.HouseRepositoryObject{}).(respository.HouseRepository),
		),
	)
}

type HouseService interface {
	Add(house model.CreateHouseRequest) (model.HouseDto, error)
	FindById(id uuid.UUID) (model.HouseDto, error)
	FindByUserId(userId uuid.UUID) []model.HouseDto
	ExistsById(id uuid.UUID) bool
}

func (h *HouseServiceObject) Add(request model.CreateHouseRequest) (response model.HouseDto, err error) {
	if err, country := h.countriesService.FindCountryByCode(request.Country); err != nil {
		return response, err
	} else if !h.userService.ExistsById(request.UserId) {
		return response, errors.New(fmt.Sprintf("user with id %s in not exists", request.UserId))
	} else {
		entity := request.ToEntity(&country)

		if entity, err := h.repository.Create(entity); err != nil {
			return response, err
		} else {
			return entity.ToDto(), nil
		}
	}
}

func (h *HouseServiceObject) FindById(id uuid.UUID) (response model.HouseDto, err error) {
	if response, err = h.repository.FindResponseById(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response, errors.New(fmt.Sprintf("house with id %s not found", id))
		}
		return response, err
	} else {
		return response, nil
	}
}

func (h *HouseServiceObject) FindByUserId(userId uuid.UUID) []model.HouseDto {
	return h.repository.FindResponseByUserId(userId)
}

func (h *HouseServiceObject) ExistsById(id uuid.UUID) bool {
	return h.repository.ExistsById(id)
}
