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
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type ProviderHandlerTestSuite struct {
	testhelper.MockTestSuite[ProviderHandler]
	providerService *mocks.ProviderService
}

func TestProviderHandlerTestSuite(t *testing.T) {
	testingSuite := &ProviderHandlerTestSuite{}
	testingSuite.TestObjectGenerator = func() ProviderHandler {
		testingSuite.providerService = new(mocks.ProviderService)
		return NewProviderHandler(testingSuite.providerService)
	}

	suite.Run(t, testingSuite)
}

func (p *ProviderHandlerTestSuite) Test_Add() {
	request := mocks.GenerateCreateProviderRequest()

	p.providerService.On("Add", request).Return(request.ToEntity().ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/providers").
		WithMethod("POST").
		WithHandler(p.TestO.Add()).
		WithBody(request)

	content := testRequest.Verify(p.T(), http.StatusCreated)

	actualResponse := model.ProviderDto{}

	json.Unmarshal(content, &actualResponse)

	assert.Equal(p.T(), model.ProviderDto{
		Id:      actualResponse.Id,
		Name:    "Name",
		Details: "Details",
		UserId:  request.UserId,
	}, actualResponse)
}

func (p *ProviderHandlerTestSuite) Test_Add_WithInvalidRequest() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/providers").
		WithMethod("POST").
		WithHandler(p.TestO.Add())

	testRequest.Verify(p.T(), http.StatusBadRequest)
}

func (p *ProviderHandlerTestSuite) Test_Add_WithErrorFromService() {
	request := mocks.GenerateCreateProviderRequest()

	err := errors.New("error")

	p.providerService.On("Add", request).Return(model.ProviderDto{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/providers").
		WithMethod("POST").
		WithHandler(p.TestO.Add()).
		WithBody(request)

	response := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), "error\n", string(response))
}

func (p *ProviderHandlerTestSuite) Test_FindById() {
	request := mocks.GenerateCreateProviderRequest()
	expected := request.ToEntity().ToDto()

	p.providerService.On("FindById", expected.Id).Return(expected, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/providers/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindById()).
		WithVar("id", expected.Id.String())

	content := testRequest.Verify(p.T(), http.StatusOK)

	actual := model.ProviderDto{}

	err := json.Unmarshal(content, &actual)

	assert.Nil(p.T(), err)

	assert.Equal(p.T(), expected, actual)
}

func (p *ProviderHandlerTestSuite) Test_FindById_WithErrorFromService() {
	p.providerService.On("FindById", mock.Anything).Return(model.ProviderDto{}, errors.New("test"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/providers/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindById()).
		WithVar("id", uuid.New().String())

	content := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), "test\n", string(content))
}

func (p *ProviderHandlerTestSuite) Test_FindById_WithNotFoundErrorFromService() {
	p.providerService.On("FindById", mock.Anything).Return(model.ProviderDto{}, int_errors.NewErrNotFound("test"))

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/providers/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindById()).
		WithVar("id", uuid.New().String())

	content := testRequest.Verify(p.T(), http.StatusNotFound)

	assert.Equal(p.T(), "test\n", string(content))
}

func (p *ProviderHandlerTestSuite) Test_FindById_WithMissingParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/providers/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindById())

	content := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), "parameter 'id' not found\n", string(content))
}

func (p *ProviderHandlerTestSuite) Test_FindById_WithInvalidUUID() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/providers/{id}").
		WithMethod("GET").
		WithHandler(p.TestO.FindById()).
		WithVar("id", "id")

	content := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), "the id is not valid id\n", string(content))
}

func (p *ProviderHandlerTestSuite) Test_FindBy() {
	expected := []model.ProviderDto{mocks.GenerateProviderDto()}
	userId := expected[0].UserId

	p.providerService.On("FindByNameLikeAndUserId", "Name", userId, 1, 15).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/providers?userId={userId}&page=1&size=15&name=Name").
		WithMethod("GET").
		WithHandler(p.TestO.FindBy()).
		WithParameter("userId", userId.String())

	content := testRequest.Verify(p.T(), http.StatusOK)

	var actual []model.ProviderDto

	err := json.Unmarshal(content, &actual)

	assert.Nil(p.T(), err)

	assert.Equal(p.T(), expected, actual)
}

func (p *ProviderHandlerTestSuite) Test_FindBy_WithInvalidId() {
	expected := []model.ProviderDto{mocks.GenerateProviderDto()}
	userId := expected[0].UserId

	p.providerService.On("FindByNameLikeAndUserId", "Name", userId, 1, 15).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/providers?userId={userId}&page=1&size=15&name=Name").
		WithMethod("GET").
		WithHandler(p.TestO.FindBy()).
		WithParameter("userId", "id")

	content := testRequest.Verify(p.T(), http.StatusBadRequest)

	assert.Equal(p.T(), "the id is not valid UUID\n", string(content))
}

func (p *ProviderHandlerTestSuite) Test_FindBy_WithDefaultValues() {
	expected := []model.ProviderDto{mocks.GenerateProviderDto()}
	userId := expected[0].UserId

	p.providerService.On("FindByNameLikeAndUserId", "", userId, 0, 25).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/providers?userId={userId}").
		WithMethod("GET").
		WithHandler(p.TestO.FindBy()).
		WithParameter("userId", userId.String())

	_ = testRequest.Verify(p.T(), http.StatusOK)
}
