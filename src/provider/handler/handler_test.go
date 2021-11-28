package handler

import (
	"encoding/json"
	"errors"
	"github.com/VlasovArtem/hob/src/provider/mocks"
	"github.com/VlasovArtem/hob/src/provider/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

var (
	providerService *mocks.ProviderService
)

func generateHandler() ProviderHandler {
	providerService = new(mocks.ProviderService)

	return NewProviderHandler(providerService)
}

func Test_Add(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCreateProviderRequest()

	providerService.On("Add", request).Return(mocks.GenerateProviderDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	content := testRequest.Verify(t, http.StatusCreated)

	actualResponse := model.ProviderDto{}

	json.Unmarshal(content, &actualResponse)

	assert.Equal(t, model.ProviderDto{
		Id:      actualResponse.Id,
		Name:    "Name",
		Details: "Details",
	}, actualResponse)
}

func Test_Add_WithInvalidRequest(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider").
		WithMethod("POST").
		WithHandler(handler.Add())

	testRequest.Verify(t, http.StatusBadRequest)
}

func Test_Add_WithErrorFromService(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCreateProviderRequest()

	err := errors.New("error")

	providerService.On("Add", request).Return(model.ProviderDto{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider").
		WithMethod("POST").
		WithHandler(handler.Add()).
		WithBody(request)

	response := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "error\n", string(response))
}

func Test_FindById(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCreateProviderRequest()
	expected := request.ToEntity().ToDto()

	providerService.On("FindById", expected.Id).Return(expected, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", expected.Id.String())

	content := testRequest.Verify(t, http.StatusOK)

	actual := model.ProviderDto{}

	err := json.Unmarshal(content, &actual)

	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}

func Test_FindById_WithErrorFromService(t *testing.T) {
	handler := generateHandler()

	providerService.On("FindById", mock.Anything).Return(model.ProviderDto{}, errors.New("test"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", uuid.New().String())

	content := testRequest.Verify(t, http.StatusNotFound)

	assert.Equal(t, "test\n", string(content))
}

func Test_FindById_WithMissingParameter(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById())

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(content))
}

func Test_FindById_WithInvalidUUID(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/{id}").
		WithMethod("GET").
		WithHandler(handler.FindById()).
		WithVar("id", "id")

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(content))
}

func Test_FindByNameLike(t *testing.T) {
	handler := generateHandler()

	expected := []model.ProviderDto{mocks.GenerateProviderDto()}

	providerService.On("FindByNameLike", "name", 1, 15).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/name/{name}?page=1&size=15").
		WithMethod("GET").
		WithHandler(handler.FindByNameLike()).
		WithVar("name", "name")

	content := testRequest.Verify(t, http.StatusOK)

	var actual []model.ProviderDto

	err := json.Unmarshal(content, &actual)

	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}

func Test_FindByNameLike_WithDefault(t *testing.T) {
	handler := generateHandler()

	expected := []model.ProviderDto{mocks.GenerateProviderDto()}

	providerService.On("FindByNameLike", "name", 0, 25).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/name/{name}").
		WithMethod("GET").
		WithHandler(handler.FindByNameLike()).
		WithVar("name", "name")

	content := testRequest.Verify(t, http.StatusOK)

	var actual []model.ProviderDto

	err := json.Unmarshal(content, &actual)

	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}

func Test_FindByNameLike_WithEmptyResponse(t *testing.T) {
	handler := generateHandler()

	var expected []model.ProviderDto

	providerService.On("FindByNameLike", "name", 0, 25).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/name/{name}").
		WithMethod("GET").
		WithHandler(handler.FindByNameLike()).
		WithVar("name", "name")

	content := testRequest.Verify(t, http.StatusOK)

	var actual []model.ProviderDto

	err := json.Unmarshal(content, &actual)

	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}
