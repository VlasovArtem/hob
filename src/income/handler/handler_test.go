package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/income/mocks"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type IncomeHandlerTestSuite struct {
	testhelper.MockTestSuite[IncomeHandler]
	incomes *mocks.IncomeService
}

func TestIncomeHandlerTestSuite(t *testing.T) {
	testingSuite := &IncomeHandlerTestSuite{}
	testingSuite.TestObjectGenerator = func() IncomeHandler {
		testingSuite.incomes = new(mocks.IncomeService)
		return NewIncomeHandler(testingSuite.incomes)
	}

	suite.Run(t, testingSuite)
}

func (i *IncomeHandlerTestSuite) Test_Add() {
	request := mocks.GenerateCreateIncomeRequest()

	i.incomes.On("Add", request).Return(request.ToEntity().ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income").
		WithMethod("POST").
		WithHandler(i.TestO.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(i.T(), http.StatusCreated)

	actual := model.IncomeDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(i.T(), model.IncomeDto{
		Id:          actual.Id,
		Name:        "Name",
		Date:        mocks.Date,
		Description: "Description",
		Sum:         100.1,
		HouseId:     actual.HouseId,
		Groups:      []groupModel.GroupDto{},
	}, actual)
}

func (i *IncomeHandlerTestSuite) Test_Add_WithNilHouseId() {
	request := mocks.GenerateCreateIncomeRequest()
	request.HouseId = nil

	i.incomes.On("Add", request).Return(request.ToEntity().ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income").
		WithMethod("POST").
		WithHandler(i.TestO.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(i.T(), http.StatusCreated)

	actual := model.IncomeDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(i.T(), model.IncomeDto{
		Id:          actual.Id,
		Name:        "Name",
		Date:        mocks.Date,
		Description: "Description",
		Sum:         100.1,
		HouseId:     nil,
		Groups:      []groupModel.GroupDto{},
	}, actual)
}

func (i *IncomeHandlerTestSuite) Test_Add_WithInvalidRequest() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income").
		WithMethod("POST").
		WithHandler(i.TestO.Add())

	testRequest.Verify(i.T(), http.StatusBadRequest)
}

func (i *IncomeHandlerTestSuite) Test_Add_WithErrorFromService() {
	request := mocks.GenerateCreateIncomeRequest()

	err := errors.New("error")
	i.incomes.On("Add", request).Return(model.IncomeDto{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income").
		WithMethod("POST").
		WithHandler(i.TestO.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(i.T(), http.StatusBadRequest)

	assert.Equal(i.T(), "error\n", string(responseByteArray))
}

func (i *IncomeHandlerTestSuite) Test_AddBatch() {
	request := mocks.GenerateCreateIncomeBatchRequest(1)

	i.incomes.On("AddBatch", request).Return(common.MapSlice(request.Incomes, func(request model.CreateIncomeRequest) model.IncomeDto {
		return request.ToEntity().ToDto()
	}), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/batch").
		WithMethod("POST").
		WithHandler(i.TestO.AddBatch()).
		WithBody(request)

	responseByteArray := testRequest.Verify(i.T(), http.StatusCreated)

	var actual []model.IncomeDto

	err := json.Unmarshal(responseByteArray, &actual)

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), []model.IncomeDto{
		{
			Id:          actual[0].Id,
			Name:        "Income Name #0",
			Date:        mocks.Date,
			Description: "Description",
			Sum:         100.1,
			HouseId:     request.Incomes[0].HouseId,
			Groups:      []groupModel.GroupDto{},
		},
	}, actual)
}

func (i *IncomeHandlerTestSuite) Test_AddBatch_WithInvalidRequest() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/batch").
		WithMethod("POST").
		WithHandler(i.TestO.AddBatch())

	testRequest.Verify(i.T(), http.StatusBadRequest)
}

func (i *IncomeHandlerTestSuite) Test_AddBatch_WithErrorFromService() {
	request := mocks.GenerateCreateIncomeBatchRequest(1)

	errorResponse := int_errors.NewBuilder().
		WithDetail("message").
		WithMessage("error")
	err := int_errors.NewErrResponse(errorResponse)

	i.incomes.On("AddBatch", request).Return([]model.IncomeDto{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/batch").
		WithMethod("POST").
		WithHandler(i.TestO.AddBatch()).
		WithBody(request)

	actual := testRequest.Verify(i.T(), http.StatusBadRequest)

	expected, err := json.Marshal(errorResponse.Build())

	assert.Nil(i.T(), err)
	assert.Equal(i.T(), append(expected, []byte("\n")...), actual)
}

func (i *IncomeHandlerTestSuite) Test_FindById() {
	response := mocks.GenerateIncomeDto()

	i.incomes.On("FindById", response.Id).
		Return(response, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/{id}").
		WithMethod("GET").
		WithHandler(i.TestO.FindById()).
		WithVar("id", response.Id.String())

	responseByteArray := testRequest.Verify(i.T(), http.StatusOK)

	actual := model.IncomeDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(i.T(), response, actual)
}

func (i *IncomeHandlerTestSuite) Test_FindById_WithError() {
	id := uuid.New()

	expected := errors.New("error")

	i.incomes.On("FindById", id).
		Return(model.IncomeDto{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/{id}").
		WithMethod("GET").
		WithHandler(i.TestO.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(i.T(), http.StatusBadRequest)

	assert.Equal(i.T(), fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func (i *IncomeHandlerTestSuite) Test_FindById_WithInvalidParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/{id}").
		WithMethod("GET").
		WithHandler(i.TestO.FindById()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(i.T(), http.StatusBadRequest)

	assert.Equal(i.T(), "the id is not valid id\n", string(responseByteArray))
}

func (i *IncomeHandlerTestSuite) Test_FindByHouseId() {
	response := []model.IncomeDto{mocks.GenerateIncomeDto()}

	i.incomes.On("FindByHouseId", *response[0].HouseId).
		Return(response, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(i.TestO.FindByHouseId()).
		WithVar("id", response[0].HouseId.String())

	responseByteArray := testRequest.Verify(i.T(), http.StatusOK)

	var actual []model.IncomeDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(i.T(), response, actual)
}

func (i *IncomeHandlerTestSuite) Test_FindByHouseId_WithEmptyResult() {
	id := uuid.New()

	i.incomes.On("FindByHouseId", id).
		Return([]model.IncomeDto{})

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(i.TestO.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(i.T(), http.StatusOK)

	var actual []model.IncomeDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(i.T(), []model.IncomeDto{}, actual)
}

func (i *IncomeHandlerTestSuite) Test_FindByHouseId_WithInvalidParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(i.TestO.FindByHouseId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(i.T(), http.StatusBadRequest)

	assert.Equal(i.T(), "the id is not valid id\n", string(responseByteArray))
}
