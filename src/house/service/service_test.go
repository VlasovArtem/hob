package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	countries "github.com/VlasovArtem/hob/src/country/service"
	groupMocks "github.com/VlasovArtem/hob/src/group/mocks"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/house/mocks"
	"github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type HouseServiceTestSuite struct {
	testhelper.MockTestSuite[HouseService]
	users            *userMocks.UserService
	houseRepository  *mocks.HouseRepository
	countriesService countries.CountryService
	groups           *groupMocks.GroupService
}

func TestHouseServiceTestSuite(t *testing.T) {
	ts := &HouseServiceTestSuite{
		countriesService: testhelper.InitCountryService(),
	}
	ts.TestObjectGenerator = func() HouseService {
		ts.users = new(userMocks.UserService)
		ts.houseRepository = new(mocks.HouseRepository)
		ts.groups = new(groupMocks.GroupService)
		return NewHouseService(ts.countriesService, ts.users, ts.houseRepository, ts.groups)
	}

	suite.Run(t, ts)
}

func (h *HouseServiceTestSuite) Test_Add() {
	request := mocks.GenerateCreateHouseRequest()

	h.users.On("ExistsById", request.UserId).Return(true)
	h.houseRepository.On("Create", mock.Anything).Return(
		func(house model.House) model.House { return house },
		nil,
	)

	actual, err := h.TestO.Add(request)

	assert.Nil(h.T(), err)
	assert.Equal(h.T(), model.HouseDto{
		Id:          actual.Id,
		Name:        "Test House",
		CountryCode: "UA",
		City:        "City",
		StreetLine1: "StreetLine1",
		StreetLine2: "StreetLine2",
		UserId:      request.UserId,
		Groups:      []groupModel.GroupDto{},
	}, actual)
}

func (h *HouseServiceTestSuite) Test_Add_WithGroupsNotFound() {
	request := mocks.GenerateCreateHouseRequest()
	request.GroupIds = []uuid.UUID{uuid.New()}

	h.users.On("ExistsById", request.UserId).Return(true)
	h.groups.On("ExistsByIds", mock.Anything).Return(false)

	income, err := h.TestO.Add(request)

	assert.Equal(h.T(), int_errors.NewErrNotFound("not all group with ids %s found", common.Join(request.GroupIds, ",")), err)
	assert.Equal(h.T(), model.HouseDto{}, income)

	h.houseRepository.AssertNotCalled(h.T(), "Create", mock.Anything)
}

func (h *HouseServiceTestSuite) Test_Add_WithUserNotFound() {
	request := mocks.GenerateCreateHouseRequest()

	h.users.On("ExistsById", request.UserId).Return(false)
	h.groups.On("ExistsByIds", mock.Anything).Return(true)

	income, err := h.TestO.Add(request)

	assert.Equal(h.T(), int_errors.NewErrNotFound("user with id %s not found", request.UserId), err)
	assert.Equal(h.T(), model.HouseDto{}, income)

	h.houseRepository.AssertNotCalled(h.T(), "Create", mock.Anything)
}

func (h *HouseServiceTestSuite) Test_AddBatch() {
	request := mocks.GenerateCreateHouseBatchRequest(2)

	h.users.On("ExistsById", mock.Anything).Return(true)
	h.houseRepository.On("CreateBatch", mock.Anything).Return(
		func(house []model.House) []model.House { return house },
		nil,
	)

	actual, err := h.TestO.AddBatch(request)

	assert.Nil(h.T(), err)
	assert.Equal(h.T(),
		[]model.HouseDto{
			{
				Id:          actual[0].Id,
				Name:        "House Name #0",
				CountryCode: "UA",
				City:        "City",
				StreetLine1: "StreetLine1",
				StreetLine2: "StreetLine2",
				UserId:      request.Houses[0].UserId,
				Groups:      []groupModel.GroupDto{},
			},
			{
				Id:          actual[1].Id,
				Name:        "House Name #1",
				CountryCode: "UA",
				City:        "City",
				StreetLine1: "StreetLine1",
				StreetLine2: "StreetLine2",
				UserId:      request.Houses[1].UserId,
				Groups:      []groupModel.GroupDto{},
			},
		}, actual)
}

func (h *HouseServiceTestSuite) Test_AddBatch_WithIssues() {
	request := mocks.GenerateCreateHouseBatchRequest(2)
	request.Houses[0].GroupIds = []uuid.UUID{uuid.New()}

	h.users.On("ExistsById", request.Houses[0].UserId).Return(false)
	h.users.On("ExistsById", request.Houses[1].UserId).Return(true)
	h.groups.On("ExistsByIds", mock.Anything).Return(false)

	actual, err := h.TestO.AddBatch(request)

	var expectedResult []model.HouseDto

	assert.Equal(h.T(), expectedResult, actual)

	builder := int_errors.NewBuilder()
	builder.WithMessage("Create house batch failed")
	builder.WithDetail(fmt.Sprintf("user with id %s not found", request.Houses[0].UserId.String()))
	builder.WithDetail(fmt.Sprintf("not all group with ids %s found", common.Join(append(request.Houses[0].GroupIds, request.Houses[1].GroupIds...), ",")))

	expectedError := int_errors.NewErrResponse(builder).(*int_errors.ErrResponse)
	actualError := err.(*int_errors.ErrResponse)

	assert.Equal(h.T(), *expectedError, *actualError)

	h.houseRepository.AssertNotCalled(h.T(), "CreateBatch", mock.Anything)
}

func (h *HouseServiceTestSuite) Test_FindById() {
	house := mocks.GenerateHouse(uuid.New())

	h.houseRepository.On("FindById", house.Id).Return(house, nil)

	actual, err := h.TestO.FindById(house.Id)

	assert.Nil(h.T(), err)
	assert.Equal(h.T(), house.ToDto(), actual)
}

func (h *HouseServiceTestSuite) Test_FindById_WithRecordNotFound() {
	id := uuid.New()

	h.houseRepository.On("FindById", id).Return(model.House{}, gorm.ErrRecordNotFound)

	actual, err := h.TestO.FindById(id)

	assert.Equal(h.T(), int_errors.NewErrNotFound("house with id %s not found", id), err)
	assert.Equal(h.T(), model.HouseDto{}, actual)
}

func (h *HouseServiceTestSuite) Test_FindById_WithRecordNotFoundExists() {
	id := uuid.New()

	expectedError := errors.New("error")
	h.houseRepository.On("FindById", id).Return(model.House{}, expectedError)

	actual, err := h.TestO.FindById(id)

	assert.Equal(h.T(), expectedError, err)
	assert.Equal(h.T(), model.HouseDto{}, actual)
}

func (h *HouseServiceTestSuite) Test_FindByUserId() {
	house := mocks.GenerateHouse(uuid.New())

	h.houseRepository.On("FindByUserId", house.UserId).Return([]model.House{house})

	actual := h.TestO.FindByUserId(house.UserId)

	assert.Equal(h.T(), []model.HouseDto{house.ToDto()}, actual)
}

func (h *HouseServiceTestSuite) Test_FindByUserId_WithNotExists() {
	id := uuid.New()

	h.houseRepository.On("FindByUserId", id).Return([]model.House{})

	actual := h.TestO.FindByUserId(id)

	assert.Equal(h.T(), []model.HouseDto{}, actual)
}

func (h *HouseServiceTestSuite) Test_ExistsById() {
	houseId := uuid.New()

	h.houseRepository.On("ExistsById", houseId).Return(true)

	assert.True(h.T(), h.TestO.ExistsById(houseId))
}

func (h *HouseServiceTestSuite) Test_ExistsById_WithNotExists() {
	id := uuid.New()

	h.houseRepository.On("ExistsById", id).Return(false)

	assert.False(h.T(), h.TestO.ExistsById(id))
}

func (h *HouseServiceTestSuite) Test_DeleteById() {
	id := uuid.New()

	h.houseRepository.On("ExistsById", id).Return(true)
	h.houseRepository.On("DeleteById", id).Return(nil)

	assert.Nil(h.T(), h.TestO.DeleteById(id))
}

func (h *HouseServiceTestSuite) Test_DeleteById_WithNotExists() {
	id := uuid.New()

	h.houseRepository.On("ExistsById", id).Return(false)

	assert.Equal(h.T(), int_errors.NewErrNotFound("house with id %s not found", id), h.TestO.DeleteById(id))

	h.houseRepository.AssertNotCalled(h.T(), "DeleteById", id)
}

func (h *HouseServiceTestSuite) Test_Update() {
	id, request := mocks.GenerateUpdateHouseRequest()

	h.houseRepository.On("ExistsById", id).Return(true)
	h.houseRepository.On("Update", id, request).Return(nil)

	assert.Nil(h.T(), h.TestO.Update(id, request))
}

func (h *HouseServiceTestSuite) Test_Update_WithErrorFromDatabase() {
	id, request := mocks.GenerateUpdateHouseRequest()

	h.houseRepository.On("ExistsById", id).Return(true)
	h.houseRepository.On("Update", id, request).Return(errors.New("test"))

	err := h.TestO.Update(id, request)
	assert.Equal(h.T(), errors.New("test"), err)
}

func (h *HouseServiceTestSuite) Test_Update_WithNotExists() {
	id, request := mocks.GenerateUpdateHouseRequest()

	h.houseRepository.On("ExistsById", id).Return(false)

	err := h.TestO.Update(id, request)
	assert.Equal(h.T(), int_errors.NewErrNotFound("house with id %s not found", id), err)

	h.houseRepository.AssertNotCalled(h.T(), "Update", id, request)
}

func (h *HouseServiceTestSuite) Test_Update_WithNotMatchingCountry() {
	id, request := mocks.GenerateUpdateHouseRequest()
	request.CountryCode = "invalid"

	h.houseRepository.On("ExistsById", id).Return(true)

	err := h.TestO.Update(id, request)
	assert.Equal(h.T(), int_errors.NewErrNotFound("country with code %s is not found", request.CountryCode), err)

	h.houseRepository.AssertNotCalled(h.T(), "Update", id, request)
}

func (h *HouseServiceTestSuite) Test_Update_WithGroupsIdsNotFound() {
	id, request := mocks.GenerateUpdateHouseRequest()
	request.GroupIds = []uuid.UUID{uuid.New()}

	h.houseRepository.On("ExistsById", id).Return(true)
	h.groups.On("ExistsByIds", mock.Anything).Return(false)

	err := h.TestO.Update(id, request)

	assert.Equal(h.T(), int_errors.NewErrNotFound("not all group with ids %s found", common.Join(request.GroupIds, ",")), err)

	h.houseRepository.AssertNotCalled(h.T(), "Update", id, request)
}
