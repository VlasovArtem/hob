package service

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	countryModel "github.com/VlasovArtem/hob/src/country/model"
	countries "github.com/VlasovArtem/hob/src/country/service"
	groupService "github.com/VlasovArtem/hob/src/group/service"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/house/repository"
	userService "github.com/VlasovArtem/hob/src/user/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

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
		dependency.FindRequiredDependency[repository.houseRepositoryStruct, repository.HouseRepository](factory),
		dependency.FindRequiredDependency[groupService.GroupServiceObject, groupService.GroupService](factory),
	)
}

type HouseService interface {
	Add(house model.CreateHouseRequest) (model.HouseDto, error)
	AddBatch(house model.CreateHouseBatchRequest) ([]model.HouseDto, error)
	FindById(id uuid.UUID) (model.HouseDto, error)
	FindByUserId(userId uuid.UUID) []model.HouseDto
	ExistsById(id uuid.UUID) bool
	DeleteById(id uuid.UUID) error
	Update(id uuid.UUID, request model.UpdateHouseRequest) error
	FindHousesByGroupId(groupId uuid.UUID) []model.HouseDto
}

func (h *HouseServiceObject) Add(request model.CreateHouseRequest) (response model.HouseDto, err error) {
	if country, err := h.countriesService.FindCountryByCode(request.CountryCode); err != nil {
		return response, err
	} else if !h.userService.ExistsById(request.UserId) {
		return response, int_errors.NewErrNotFound("user with id %s not found", request.UserId)
	} else if len(request.GroupIds) != 0 && !h.groupService.ExistsByIds(request.GroupIds) {
		return response, int_errors.NewErrNotFound("not all group with ids %s found", common.Join(request.GroupIds, ","))
	} else {
		entity := request.ToEntity(&country)

		if entity, err := h.houseRepository.Create(entity, "Groups.*"); err != nil {
			return response, err
		} else {
			return entity.ToDto(), nil
		}
	}
}

func (h *HouseServiceObject) AddBatch(request model.CreateHouseBatchRequest) ([]model.HouseDto, error) {
	if len(request.Houses) == 0 {
		return make([]model.HouseDto, 0), nil
	}

	countryShortName := make(map[string]bool)
	userIds := make(map[uuid.UUID]bool)
	groups := make(map[uuid.UUID]bool)

	for _, createHouseRequest := range request.Houses {
		countryShortName[createHouseRequest.CountryCode] = true
		userIds[createHouseRequest.UserId] = true
		if len(createHouseRequest.GroupIds) != 0 {
			for _, groupId := range createHouseRequest.GroupIds {
				groups[groupId] = true
			}
		}
	}

	builder := int_errors.NewBuilder()

	countryResult := make(map[string]countryModel.Country)

	for code, _ := range countryShortName {
		if countryByCode, err := h.countriesService.FindCountryByCode(code); err != nil {
			builder.WithDetail(err.Error())
		} else {
			countryResult[code] = countryByCode
		}
	}

	for userId, _ := range userIds {
		if !h.userService.ExistsById(userId) {
			builder.WithDetail(fmt.Sprintf("user with id %s not found", userId))
		}
	}

	var groupIds []uuid.UUID

	for groupId := range groups {
		groupIds = append(groupIds, groupId)
	}

	if len(groupIds) != 0 && !h.groupService.ExistsByIds(groupIds) {
		builder.WithDetail(fmt.Sprintf("not all group with ids %s found", common.Join(groupIds, ",")))
	}

	if builder.HasErrors() {
		return nil, int_errors.NewErrResponse(builder.WithMessage("Create house batch failed"))
	}

	entitiesForCreation := common.MapSlice(request.Houses, func(request model.CreateHouseRequest) model.House {
		if country, ok := countryResult[request.CountryCode]; !ok {
			log.Error().Msg(fmt.Sprintf("Country with code %s not found", request.CountryCode))
			return model.House{}
		} else {
			return request.ToEntity(&country)
		}
	})

	if createBatch, err := h.houseRepository.CreateBatch(entitiesForCreation); err != nil {
		return nil, err
	} else {
		return common.MapSlice(createBatch, func(entity model.House) model.HouseDto {
			return entity.ToDto()
		}), nil
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

	return common.MapSlice(houseEntities, func(entity model.House) model.HouseDto {
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

func (h *HouseServiceObject) FindHousesByGroupId(groupId uuid.UUID) []model.HouseDto {
	return common.MapSlice(h.houseRepository.FindHousesByGroupId(groupId), func(entity model.House) model.HouseDto {
		return entity.ToDto()
	})
}
