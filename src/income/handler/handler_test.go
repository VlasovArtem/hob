package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/income/mocks"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	incomes *mocks.IncomeService
)

func handlerGenerator() IncomeHandler {
	incomes = new(mocks.IncomeService)

	return NewIncomeHandler(incomes)
}

func Test_AddIncome(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateCreateIncomeRequest()

	incomes.On("Add", request).Return(request.ToEntity().ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusCreated)

	actual := model.IncomeDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, model.IncomeDto{
		Id:          actual.Id,
		Name:        "Name",
		Date:        mocks.Date,
		Description: "Description",
		Sum:         100.1,
		HouseId:     actual.HouseId,
		Groups:      []groupModel.GroupDto{},
	}, actual)
}

func Test_AddIncome_WithInvalidRequest(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_AddIncome_WithErrorFromService(t *testing.T) {
	handler := handlerGenerator()

	request := mocks.GenerateCreateIncomeRequest()

	err := errors.New("error")
	incomes.On("Add", request).Return(model.IncomeDto{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "error\n", string(responseByteArray))
}

func Test_FindById(t *testing.T) {
	handler := handlerGenerator()

	response := mocks.GenerateIncomeDto()

	incomes.On("FindById", response.Id).
		Return(response, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", response.Id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	actual := model.IncomeDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, response, actual)
}

func Test_FindById_WithError(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	expected := errors.New("error")

	incomes.On("FindById", id).
		Return(model.IncomeDto{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func Test_FindById_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/income/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}

func Test_FindByHouseId(t *testing.T) {
	handler := handlerGenerator()

	response := []model.IncomeDto{mocks.GenerateIncomeDto()}

	incomes.On("FindByHouseId", response[0].HouseId).
		Return(response, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", response[0].HouseId.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []model.IncomeDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, response, actual)
}

func Test_FindByHouseId_WithEmptyResult(t *testing.T) {
	handler := handlerGenerator()

	id := uuid.New()

	incomes.On("FindByHouseId", id).
		Return([]model.IncomeDto{})

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(t, http.StatusOK)

	var actual []model.IncomeDto

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(t, []model.IncomeDto{}, actual)
}

func Test_FindByHouseId_WithInvalidParameter(t *testing.T) {
	handler := handlerGenerator()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meter/house/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByHouseId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(responseByteArray))
}
