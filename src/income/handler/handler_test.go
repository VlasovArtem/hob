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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

var nilTime *time.Time

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
		WithURL("https://test.com/api/v1/incomes").
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
		WithURL("https://test.com/api/v1/incomes").
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
		WithURL("https://test.com/api/v1/incomes").
		WithMethod("POST").
		WithHandler(i.TestO.Add())

	testRequest.Verify(i.T(), http.StatusBadRequest)
}

func (i *IncomeHandlerTestSuite) Test_Add_WithErrorFromService() {
	request := mocks.GenerateCreateIncomeRequest()

	err := errors.New("error")
	i.incomes.On("Add", request).Return(model.IncomeDto{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes").
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
		WithURL("https://test.com/api/v1/incomes/batch").
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
		WithURL("https://test.com/api/v1/incomes/batch").
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
		WithURL("https://test.com/api/v1/incomes/batch").
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
		WithURL("https://test.com/api/v1/incomes/{id}").
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
		WithURL("https://test.com/api/v1/incomes/{id}").
		WithMethod("GET").
		WithHandler(i.TestO.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(i.T(), http.StatusBadRequest)

	assert.Equal(i.T(), fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func (i *IncomeHandlerTestSuite) Test_FindById_WithInvalidParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes/{id}").
		WithMethod("GET").
		WithHandler(i.TestO.FindById()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(i.T(), http.StatusBadRequest)

	assert.Equal(i.T(), "the id is not valid id\n", string(responseByteArray))
}

func (i *IncomeHandlerTestSuite) Test_FindByHouseId() {
	response := []model.IncomeDto{mocks.GenerateIncomeDto()}
	from, fromString, to, toString := createFromAndTo()

	i.incomes.On("FindByHouseId", *response[0].HouseId, 10, 0, from, to).
		Return(response, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes/house/{id}?limit={limit}&offset={offset}&from={from}&to={to}").
		WithMethod("GET").
		WithHandler(i.TestO.FindByHouseId()).
		WithVar("id", response[0].HouseId.String()).
		WithParameter("limit", "10").
		WithParameter("offset", "0").
		WithParameter("from", fromString).
		WithParameter("to", toString)

	responseByteArray := testRequest.Verify(i.T(), http.StatusOK)

	var actual []model.IncomeDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(i.T(), response, actual)
}

func (i *IncomeHandlerTestSuite) Test_FindByHouseId_WithFrom() {
	response := []model.IncomeDto{mocks.GenerateIncomeDto()}
	from, fromString, _, _ := createFromAndTo()

	i.incomes.On("FindByHouseId", *response[0].HouseId, 10, 0, from, nilTime).
		Return(response, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes/house/{id}?limit={limit}&offset={offset}&from={from}").
		WithMethod("GET").
		WithHandler(i.TestO.FindByHouseId()).
		WithVar("id", response[0].HouseId.String()).
		WithParameter("limit", "10").
		WithParameter("offset", "0").
		WithParameter("from", fromString)

	responseByteArray := testRequest.Verify(i.T(), http.StatusOK)

	var actual []model.IncomeDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(i.T(), response, actual)
}

func (i *IncomeHandlerTestSuite) Test_FindByHouseId_WithDefaultLimitAndOffset() {
	response := []model.IncomeDto{mocks.GenerateIncomeDto()}

	i.incomes.On("FindByHouseId", *response[0].HouseId, 25, 0, nilTime, nilTime).
		Return(response, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes/house/{id}").
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

	i.incomes.On("FindByHouseId", id, 25, 0, nilTime, nilTime).
		Return([]model.IncomeDto{})

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes/house/{id}").
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
		WithURL("https://test.com/api/v1/incomes/house/{id}").
		WithMethod("GET").
		WithHandler(i.TestO.FindByHouseId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(i.T(), http.StatusBadRequest)

	assert.Equal(i.T(), "the id is not valid id\n", string(responseByteArray))
}

func (i *IncomeHandlerTestSuite) Test_Update() {
	id, request := mocks.GenerateUpdateIncomeRequest()

	i.incomes.On("Update", id, request).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes/{id}").
		WithMethod("PUT").
		WithHandler(i.TestO.Update()).
		WithBody(request).
		WithVar("id", id.String())

	testRequest.Verify(i.T(), http.StatusOK)
}

func (i *IncomeHandlerTestSuite) Test_Update_WithInvalidId() {
	_, request := mocks.GenerateUpdateIncomeRequest()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes/{id}").
		WithMethod("PUT").
		WithHandler(i.TestO.Update()).
		WithBody(request).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(i.T(), http.StatusBadRequest)

	assert.Equal(i.T(), "the id is not valid id\n", string(responseByteArray))

	i.incomes.AssertNotCalled(i.T(), "Update", mock.Anything, mock.Anything)
}

func (i *IncomeHandlerTestSuite) Test_Update_WithInvalidRequest() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes/{id}").
		WithMethod("PUT").
		WithHandler(i.TestO.Update()).
		WithVar("id", uuid.New().String())

	testRequest.Verify(i.T(), http.StatusBadRequest)
}

func (i *IncomeHandlerTestSuite) Test_Update_WithErrorFromService() {
	id, request := mocks.GenerateUpdateIncomeRequest()

	expected := errors.New("error")

	i.incomes.On("Update", id, request).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes/{id}").
		WithMethod("PUT").
		WithHandler(i.TestO.Update()).
		WithVar("id", id.String()).
		WithBody(request)

	responseByteArray := testRequest.Verify(i.T(), http.StatusBadRequest)

	assert.Equal(i.T(), fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func (i *IncomeHandlerTestSuite) Test_Delete() {
	id := uuid.New()

	i.incomes.On("DeleteById", id).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes/{id}").
		WithMethod("DELETE").
		WithHandler(i.TestO.Delete()).
		WithVar("id", id.String())

	testRequest.Verify(i.T(), http.StatusNoContent)
}

func (i *IncomeHandlerTestSuite) Test_Delete_WithMissingParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes/{id}").
		WithMethod("DELETE").
		WithHandler(i.TestO.Delete())

	responseByteArray := testRequest.Verify(i.T(), http.StatusBadRequest)

	assert.Equal(i.T(), "parameter 'id' not found\n", string(responseByteArray))
}

func (i *IncomeHandlerTestSuite) Test_Delete_WithInvalidParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/incomes/{id}").
		WithMethod("DELETE").
		WithHandler(i.TestO.Delete()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(i.T(), http.StatusBadRequest)

	assert.Equal(i.T(), "the id is not valid id\n", string(responseByteArray))
}

func createFromAndTo() (from *time.Time, fromString string, to *time.Time, toString string) {
	fromString = time.Now().UTC().Format(time.RFC3339)
	toString = time.Now().UTC().Format(time.RFC3339)
	fromDate, _ := time.Parse(time.RFC3339, fromString)
	toDate, _ := time.Parse(time.RFC3339, toString)
	return &fromDate, fromString, &toDate, toString
}
