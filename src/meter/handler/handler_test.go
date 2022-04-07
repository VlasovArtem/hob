package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/meter/mocks"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type MeterHandlerTestSuite struct {
	testhelper.MockTestSuite[MeterHandler]
	meters *mocks.MeterService
}

func TestMeterHandlerTestSuite(t *testing.T) {
	testingSuite := &MeterHandlerTestSuite{}
	testingSuite.TestObjectGenerator = func() MeterHandler {
		testingSuite.meters = new(mocks.MeterService)
		return NewMeterHandler(testingSuite.meters)
	}

	suite.Run(t, testingSuite)
}

func (m *MeterHandlerTestSuite) Test_AddMeter() {
	request := mocks.GenerateCreateMeterRequest()

	m.meters.On("Add", request).Return(request.ToEntity().ToDto(), nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters").
		WithMethod("POST").
		WithHandler(m.TestO.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(m.T(), http.StatusCreated)

	actual := model.MeterDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(m.T(), model.MeterDto{
		Id:   actual.Id,
		Name: "Name",
		Details: map[string]float64{
			"first":  1.1,
			"second": 2.2,
		},
		Description: "Description",
		PaymentId:   request.PaymentId,
	}, actual)
}

func (m *MeterHandlerTestSuite) Test_AddMeter_WithInvalidRequest() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters").
		WithMethod("POST").
		WithHandler(m.TestO.Add())

	testRequest.Verify(m.T(), http.StatusBadRequest)
}

func (m *MeterHandlerTestSuite) Test_AddMeter_WithErrorFromService() {
	request := mocks.GenerateCreateMeterRequest()

	err := errors.New("error")
	m.meters.On("Add", request).Return(model.MeterDto{}, err)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters").
		WithMethod("POST").
		WithHandler(m.TestO.Add()).
		WithBody(request)

	responseByteArray := testRequest.Verify(m.T(), http.StatusBadRequest)

	assert.Equal(m.T(), "error\n", string(responseByteArray))
}

func (m *MeterHandlerTestSuite) Test_FindById() {
	id := uuid.New()

	meterResponse := mocks.GenerateMeterResponse(id)

	m.meters.On("FindById", id).
		Return(meterResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/{id}").
		WithMethod("GET").
		WithHandler(m.TestO.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(m.T(), http.StatusOK)

	actual := model.MeterDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(m.T(), meterResponse, actual)
}

func (m *MeterHandlerTestSuite) Test_FindById_WithError() {
	id := uuid.New()

	expected := errors.New("error")

	m.meters.On("FindById", id).
		Return(model.MeterDto{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/{id}").
		WithMethod("GET").
		WithHandler(m.TestO.FindById()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(m.T(), http.StatusBadRequest)

	assert.Equal(m.T(), fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func (m *MeterHandlerTestSuite) Test_FindById_WithInvalidParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/{id}").
		WithMethod("GET").
		WithHandler(m.TestO.FindById()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(m.T(), http.StatusBadRequest)

	assert.Equal(m.T(), "the id is not valid id\n", string(responseByteArray))
}

func (m *MeterHandlerTestSuite) Test_FindByPaymentId() {
	id := uuid.New()

	meterResponse := mocks.GenerateMeterResponse(id)

	m.meters.On("FindByPaymentId", id).
		Return(meterResponse, nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/payment/{id}").
		WithMethod("GET").
		WithHandler(m.TestO.FindByPaymentId()).
		WithVar("id", id.String())

	responseByteArray := testRequest.Verify(m.T(), http.StatusOK)

	actual := model.MeterDto{}

	json.Unmarshal(responseByteArray, &actual)

	assert.Equal(m.T(), meterResponse, actual)
}

func (m *MeterHandlerTestSuite) Test_FindByPaymentId_WithError() {
	id := uuid.New()

	expected := errors.New("error")

	m.meters.On("FindByPaymentId", id).
		Return(model.MeterDto{}, expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/payment/{id}").
		WithMethod("GET").
		WithHandler(m.TestO.FindByPaymentId()).
		WithVar("id", id.String())

	body := testRequest.Verify(m.T(), http.StatusBadRequest)

	assert.Equal(m.T(), fmt.Sprintf("%s\n", expected.Error()), string(body))
}

func (m *MeterHandlerTestSuite) Test_FindByPaymentId_WithInvalidParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/payment/{id}").
		WithMethod("GET").
		WithHandler(m.TestO.FindByPaymentId()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(m.T(), http.StatusBadRequest)

	assert.Equal(m.T(), "the id is not valid id\n", string(responseByteArray))
}

func (m *MeterHandlerTestSuite) Test_Update() {
	id, request := mocks.GenerateUpdateMeterRequest()

	m.meters.On("Update", id, request).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/{id}").
		WithMethod("PUT").
		WithHandler(m.TestO.Update()).
		WithBody(request).
		WithVar("id", id.String())

	testRequest.Verify(m.T(), http.StatusOK)
}

func (m *MeterHandlerTestSuite) Test_Update_WithInvalidId() {
	_, request := mocks.GenerateUpdateMeterRequest()

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/{id}").
		WithMethod("PUT").
		WithHandler(m.TestO.Update()).
		WithBody(request).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(m.T(), http.StatusBadRequest)

	assert.Equal(m.T(), "the id is not valid id\n", string(responseByteArray))

	m.meters.AssertNotCalled(m.T(), "Update", mock.Anything, mock.Anything)
}

func (m *MeterHandlerTestSuite) Test_Update_WithInvalidRequest() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/{id}").
		WithMethod("PUT").
		WithHandler(m.TestO.Update()).
		WithVar("id", uuid.New().String())

	testRequest.Verify(m.T(), http.StatusBadRequest)
}

func (m *MeterHandlerTestSuite) Test_Update_WithErrorFromService() {
	id, request := mocks.GenerateUpdateMeterRequest()

	expected := errors.New("error")

	m.meters.On("Update", id, request).Return(expected)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/{id}").
		WithMethod("PUT").
		WithHandler(m.TestO.Update()).
		WithVar("id", id.String()).
		WithBody(request)

	responseByteArray := testRequest.Verify(m.T(), http.StatusBadRequest)

	assert.Equal(m.T(), fmt.Sprintf("%s\n", expected.Error()), string(responseByteArray))
}

func (m *MeterHandlerTestSuite) Test_Delete() {
	id := uuid.New()

	m.meters.On("DeleteById", id).Return(nil)

	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/{id}").
		WithMethod("DELETE").
		WithHandler(m.TestO.Delete()).
		WithVar("id", id.String())

	testRequest.Verify(m.T(), http.StatusNoContent)
}

func (m *MeterHandlerTestSuite) Test_Delete_WithMissingParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/{id}").
		WithMethod("DELETE").
		WithHandler(m.TestO.Delete())

	responseByteArray := testRequest.Verify(m.T(), http.StatusBadRequest)

	assert.Equal(m.T(), "parameter 'id' not found\n", string(responseByteArray))
}

func (m *MeterHandlerTestSuite) Test_Delete_WithInvalidParameter() {
	testRequest := testhelper.NewTestRequest().
		WithURL("https://test.com/api/v1/meters/{id}").
		WithMethod("DELETE").
		WithHandler(m.TestO.Delete()).
		WithVar("id", "id")

	responseByteArray := testRequest.Verify(m.T(), http.StatusBadRequest)

	assert.Equal(m.T(), "the id is not valid id\n", string(responseByteArray))
}
