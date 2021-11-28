package handler

import (
	"encoding/json"
	"errors"
	"github.com/VlasovArtem/hob/src/provider/custom/mocks"
	"github.com/VlasovArtem/hob/src/provider/custom/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

var (
	providerService *mocks.CustomProviderService
)

func generateHandler() CustomProviderHandler {
	providerService = new(mocks.CustomProviderService)

	return NewCustomProviderHandler(providerService)
}

func Test_Add(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCustomProviderRequest()

	providerService.On("Add", request).Return(mocks.GenerateCustomProviderDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/custom").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	content := testRequest.Verify(t, http.StatusCreated)

	actualResponse := model.CustomProviderDto{}

	json.Unmarshal(content, &actualResponse)

	assert.Equal(t, model.CustomProviderDto{
		Id:      actualResponse.Id,
		Name:    "Name",
		Details: "Details",
		UserId:  actualResponse.UserId,
	}, actualResponse)
}

func Test_Add_WithInvalidRequest(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/custom").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_Add_WithErrorFromService(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCustomProviderRequest()

	err := errors.New("error")

	providerService.On("Add", request).Return(model.CustomProviderDto{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/custom").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	response := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "error\n", string(response))
}

func Test_FindById(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCustomProviderRequest()
	expected := request.ToEntity().ToDto()

	providerService.On("FindById", expected.Id).Return(expected, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/custom/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", expected.Id.String())

	content := testRequest.Verify(t, http.StatusOK)

	actual := model.CustomProviderDto{}

	err := json.Unmarshal(content, &actual)

	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}

func Test_FindById_WithErrorFromService(t *testing.T) {
	handler := generateHandler()

	providerService.On("FindById", mock.Anything).Return(model.CustomProviderDto{}, errors.New("test"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/custom/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", uuid.New().String())

	content := testRequest.Verify(t, http.StatusNotFound)

	assert.Equal(t, "test\n", string(content))
}

func Test_FindById_WithMissingParameter(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/custom/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById())

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(content))
}

func Test_FindById_WithInvalidUUID(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/custom/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(content))
}

func Test_FindByUserId(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCustomProviderRequest()
	expected := []model.CustomProviderDto{request.ToEntity().ToDto()}

	providerService.On("FindByUserId", expected[0].UserId).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/custom/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", expected[0].UserId.String())

	content := testRequest.Verify(t, http.StatusOK)

	var actual []model.CustomProviderDto

	err := json.Unmarshal(content, &actual)

	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}

func Test_FindByUserId_WithEmptyResponse(t *testing.T) {
	handler := generateHandler()

	providerService.On("FindByUserId", mock.Anything).Return([]model.CustomProviderDto{})

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/custom/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", uuid.New().String())

	content := testRequest.Verify(t, http.StatusOK)

	var actual []model.CustomProviderDto

	err := json.Unmarshal(content, &actual)

	assert.Nil(t, err)

	assert.Equal(t, []model.CustomProviderDto{}, actual)
}

func Test_FindByUserId_WithMissingParameter(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/custom/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById())

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(content))
}

func Test_FindByUserId_WithInvalidUUID(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/custom/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(content))
}
