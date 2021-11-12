package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"payment/model"
	"test/mock"
	"test/testhelper"
	"testing"
	"time"
)

var (
	users   *mock.UserServiceMock
	houses  *mock.HouseServiceMock
	houseId = testhelper.ParseUUID("d077adaa-00d7-4e80-ac86-57512267505d")
	userId  = testhelper.ParseUUID("ad2c5035-6745-48d0-9eee-fd22f5dae8e0")
	date    = time.Now()
)

func serviceGenerator() PaymentService {
	users = new(mock.UserServiceMock)
	houses = new(mock.HouseServiceMock)

	return NewPaymentService(users, houses)
}

func Test_AddPayment(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", userId).Return(true)
	houses.On("ExistsById", houseId).Return(true)

	request := generateCreatePaymentRequest()

	payment, err := paymentService.AddPayment(request)

	expectedEntity := request.ToEntity()
	expectedEntity.Id = payment.Id
	expectedResponse := expectedEntity.ToResponse()

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, payment)
}

func Test_AddPayment_WithUserNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", userId).Return(false)

	request := generateCreatePaymentRequest()

	payment, err := paymentService.AddPayment(request)

	assert.Equal(t, errors.New(fmt.Sprintf("user with id %s in not exists", request.UserId)), err)
	assert.Equal(t, model.PaymentResponse{}, payment)
}

func Test_AddPayment_WithHouseNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", userId).Return(true)
	houses.On("ExistsById", houseId).Return(false)

	request := generateCreatePaymentRequest()

	payment, err := paymentService.AddPayment(request)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s in not exists", request.HouseId)), err)
	assert.Equal(t, model.PaymentResponse{}, payment)
}

func Test_FindPaymentById(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", userId).Return(true)
	houses.On("ExistsById", houseId).Return(true)

	request := generateCreatePaymentRequest()

	payment, err := paymentService.AddPayment(request)

	assert.Nil(t, err)

	actual, err := paymentService.FindPaymentById(payment.Id)

	assert.Nil(t, nil)
	assert.Equal(t, payment, actual)
}

func Test_FindPaymentById_WithNotExistingId(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()
	actual, err := paymentService.FindPaymentById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("payment with id %s not found", id)), err)
	assert.Equal(t, model.PaymentResponse{}, actual)
}

func Test_FindPaymentByHouseId(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", userId).Return(true)
	houses.On("ExistsById", houseId).Return(true)

	request := generateCreatePaymentRequest()

	payment, err := paymentService.AddPayment(request)

	assert.Nil(t, err)

	payments := paymentService.FindPaymentByHouseId(payment.HouseId)

	assert.Equal(t, []model.PaymentResponse{payment}, payments)
}

func Test_FindPaymentByHouseId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	payments := paymentService.FindPaymentByHouseId(uuid.New())

	assert.Equal(t, []model.PaymentResponse{}, payments)
}

func Test_FindPaymentByUserId(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", userId).Return(true)
	houses.On("ExistsById", houseId).Return(true)

	request := generateCreatePaymentRequest()

	payment, err := paymentService.AddPayment(request)

	assert.Nil(t, err)

	payments := paymentService.FindPaymentByUserId(payment.UserId)

	assert.Equal(t, []model.PaymentResponse{payment}, payments)
}

func Test_FindPaymentByUserId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	payments := paymentService.FindPaymentByUserId(uuid.New())

	assert.Equal(t, []model.PaymentResponse{}, payments)
}

func Test_ExistsById(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	paymentService.(*paymentServiceObject).payments[id] = model.Payment{}

	assert.True(t, paymentService.ExistsById(id))
}

func Test_ExistsById_WithNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	assert.False(t, paymentService.ExistsById(id))
}

func generateCreatePaymentRequest() model.CreatePaymentRequest {
	return model.CreatePaymentRequest{
		Name:        "Test Payment",
		Description: "Test Payment Description",
		HouseId:     houseId,
		UserId:      userId,
		Date:        date,
		Sum:         1000,
	}
}
