package service

import (
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	countries "github.com/VlasovArtem/hob/src/country/service"
	groupService "github.com/VlasovArtem/hob/src/group/service"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/house/repository"
	userService "github.com/VlasovArtem/hob/src/user/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"reflect"
)

var HouseServiceType = reflect.TypeOf(HouseServiceObject{})

type HouseServiceObject struct {
	countriesService countries.CountryService
	userService      userService.UserService
	houseRepository  repository.HouseRepository
	groupService     groupService.GroupService
}

func NewHouseService(
	countriesService countries.CountryService,
	userService userService.UserService,
	repository repository.HouseRepository,
	groupService groupService.GroupService,
) HouseService {
	if countriesService == nil {
		log.Fatal().Msg("CountryCode service is required")
	}

	return &HouseServiceObject{
		countriesService: countriesService,
		userService:      userService,
		houseRepository:  repository,
		groupService:     groupService,
	}
}

func (h *HouseServiceObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewHouseService(
		dependency.FindRequiredDependency[countries.CountryServiceObject, countries.CountryService](factory),
		dependency.FindRequiredDependency[userService.UserServiceObject, userService.UserService](factory),
		dependency.FindRequiredDependency[repository.HouseRepositoryObject, repository.HouseRepository](factory),
		dependency.FindRequiredDependency[groupService.GroupServiceObject, groupService.GroupService](factory),
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
	if country, err := h.countriesService.FindCountryByCode(request.CountryCode); err != nil {
		return response, err
	} else if !h.userService.ExistsById(request.UserId) {
		return response, int_errors.NewErrNotFound("user with id %s in not exists", request.UserId)
	} else if len(request.GroupIds) != 0 && !h.groupService.ExistsByIds(request.GroupIds) {
		return response, int_errors.NewErrNotFound("not all group with ids %s found", common.Join(request.GroupIds, ","))
	} else {
		entity := request.ToEntity(&country)

		if entity, err := h.houseRepository.Create(entity); err != nil {
			return response, err
		} else {
			return entity.ToDto(), nil
		}
	}
}

func (h *HouseServiceObject) FindById(id uuid.UUID) (response model.HouseDto, err error) {
	if entity, err := h.houseRepository.FindById(id); err != nil {
		return response, database.HandlerFindError(err, "house with id %s not found", id)
	} else {
		return entity.ToDto(), nil
	}
}

func (h *HouseServiceObject) FindByUserId(userId uuid.UUID) []model.HouseDto {
	houseEntities := h.houseRepository.FindByUserId(userId)

	return common.Map(houseEntities, func(entity model.House) model.HouseDto {
		return entity.ToDto()
	})
}

func (h *HouseServiceObject) ExistsById(id uuid.UUID) bool {
	return h.houseRepository.ExistsById(id)
}

func (h *HouseServiceObject) DeleteById(id uuid.UUID) error {
	if !h.ExistsById(id) {
		return int_errors.NewErrNotFound("house with id %s not found", id)
	}
	return h.houseRepository.DeleteById(id)
}

func (h *HouseServiceObject) Update(id uuid.UUID, request model.UpdateHouseRequest) error {
	if !h.ExistsById(id) {
		return int_errors.NewErrNotFound("house with id %s not found", id)
	}
	if len(request.GroupIds) != 0 && !h.groupService.ExistsByIds(request.GroupIds) {
		return int_errors.NewErrNotFound("not all group with ids %s found", common.Join(request.GroupIds, ","))
	}
	if _, err := h.countriesService.FindCountryByCode(request.CountryCode); err != nil {
		return err
	} else {
		return h.houseRepository.Update(id, request)
	}
}
