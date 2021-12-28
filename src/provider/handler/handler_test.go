package handler

import (
	"encoding/json"
	"errors"
	int_errors "github.com/VlasovArtem/hob/src/common/int-errors"
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

	providerService.On("Add", request).Return(request.ToEntity().ToDto(), nil)

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
		UserId:  request.UserId,
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

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "test\n", string(content))
}

func Test_FindById_WithNotFoundErrorFromService(t *testing.T) {
	handler := generateHandler()

	providerService.On("FindById", mock.Anything).Return(model.ProviderDto{}, int_errors.NewErrNotFound("test"))

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

func Test_FindByUserId(t *testing.T) {
	handler := generateHandler()

	request := mocks.GenerateCreateProviderRequest()
	expected := []model.ProviderDto{request.ToEntity().ToDto()}
	userId := expected[0].UserId

	providerService.On("FindByUserId", userId).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", userId.String())

	content := testRequest.Verify(t, http.StatusOK)

	var actual []model.ProviderDto

	err := json.Unmarshal(content, &actual)

	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}

func Test_FindByUserId_WithMissingParameter(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId())

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(content))
}

func Test_FindByUserId_WithInvalidUUID(t *testing.T) {
	handler := generateHandler()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/user/{id}").
		WithMethod("GET").
		WithHandler(handler.FindByUserId()).
		WithVar("id", "id")

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(content))
}

func Test_FindByUserIdAndNameLike(t *testing.T) {
	handler := generateHandler()

	expected := []model.ProviderDto{mocks.GenerateProviderDto()}
	userId := expected[0].UserId

	providerService.On("FindByNameLikeAndUserId", "Name", userId, 1, 15).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/user/{id}?page=1&size=15").
		WithMethod("POST").
		WithHandler(handler.FindByNameLikeAndUserId()).
		WithVar("id", userId.String()).
		WithBody(FindByNameRequest{"Name"})

	content := testRequest.Verify(t, http.StatusOK)

	var actual []model.ProviderDto

	err := json.Unmarshal(content, &actual)

	assert.Nil(t, err)

	assert.Equal(t, expected, actual)
}

func Test_FindByUserIdAndNameLike_WithInvalidId(t *testing.T) {
	handler := generateHandler()

	expected := []model.ProviderDto{mocks.GenerateProviderDto()}
	userId := expected[0].UserId

	providerService.On("FindByNameLikeAndUserId", "Name", userId, 1, 15).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/user/{id}?page=1&size=15").
		WithMethod("POST").
		WithHandler(handler.FindByNameLikeAndUserId()).
		WithVar("id", "id").
		WithBody(FindByNameRequest{"Name"})

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "the id is not valid id\n", string(content))
}

func Test_FindByUserIdAndNameLike_WithoutId(t *testing.T) {
	handler := generateHandler()

	expected := []model.ProviderDto{mocks.GenerateProviderDto()}
	userId := expected[0].UserId

	providerService.On("FindByNameLikeAndUserId", "Name", userId, 1, 15).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/user/{id}?page=1&size=15").
		WithMethod("POST").
		WithHandler(handler.FindByNameLikeAndUserId()).
		WithBody(FindByNameRequest{"Name"})

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "parameter 'id' not found\n", string(content))
}

func Test_FindByUserIdAndNameLike_WithoutBody(t *testing.T) {
	handler := generateHandler()

	expected := []model.ProviderDto{mocks.GenerateProviderDto()}
	userId := expected[0].UserId

	providerService.On("FindByNameLikeAndUserId", "Name", userId, 1, 15).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/provider/user/{id}?page=1&size=15").
		WithMethod("POST").
		WithHandler(handler.FindByNameLikeAndUserId()).
		WithVar("id", userId.String())

	content := testRequest.Verify(t, http.StatusBadRequest)

	assert.Equal(t, "body not found\n", string(content))
}
