package service

import (
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	countries "github.com/VlasovArtem/hob/src/country/service"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/house/respository"
	userService "github.com/VlasovArtem/hob/src/user/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"reflect"
)

var HouseServiceType = reflect.TypeOf(HouseServiceObject{})

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
		log.Fatal().Msg("Country service is required")
	}

	return &HouseServiceObject{
		countriesService: countriesService,
		userService:      userService,
		repository:       repository,
	}
}

func (h *HouseServiceObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewHouseService(
		factory.FindRequiredByObject(countries.CountryServiceObject{}).(countries.CountryService),
		factory.FindRequiredByType(userService.UserServiceType).(userService.UserService),
		factory.FindRequiredByType(respository.HouseRepositoryType).(respository.HouseRepository),
	)
}

type HouseService interface {
	Add(house model.CreateHouseRequest) (model.HouseDto, error)
	FindById(id uuid.UUID) (model.HouseDto, error)
	FindByUserId(userId uuid.UUID) []model.HouseDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateHouseRequest) error
}

func (h *HouseServiceObject) Add(request model.CreateHouseRequest) (response model.HouseDto, err error) {
	if err, country := h.countriesService.FindCountryByCode(request.Country); err != nil {
		return response, err
	} else if !h.userService.ExistsById(request.UserId) {
		return response, int_errors.NewErrNotFound("user with id %s in not exists", request.UserId)
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
	if response, err = h.repository.FindDtoById(id); err != nil {
		return response, database.HandlerFindError(err, "house with id %s not found", id)
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

func (h *HouseServiceObject) DeleteById(id uuid.UUID) error {
	if !h.ExistsById(id) {
		return int_errors.NewErrNotFound("house with id %s not found", id)
	}
	return h.repository.DeleteById(id)
}

func (h *HouseServiceObject) Update(id uuid.UUID, request model.UpdateHouseRequest) error {
	if !h.ExistsById(id) {
		return int_errors.NewErrNotFound("house with id %s not found", id)
	}
	if err, _ := h.countriesService.FindCountryByCode(request.Country); err != nil {
		return err
	} else {
		return h.repository.Update(id, request)
	}
}
