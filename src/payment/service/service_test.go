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

	payment, err := paymentService.Add(request)

	expectedEntity := request.ToEntity()
	expectedEntity.Id = payment.Id
	expectedResponse := expectedEntity.ToResponse()

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, payment)

	serviceObject := paymentService.(*paymentServiceObject)

	_, paymentExists := serviceObject.payments[payment.Id]
	assert.True(t, paymentExists)

	_, housePaymentExists := serviceObject.housePayments[payment.HouseId]
	assert.True(t, housePaymentExists)

	_, userPaymentExists := serviceObject.userPayments[payment.UserId]
	assert.True(t, userPaymentExists)
}

func Test_Add_WithUserNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", userId).Return(false)

	request := generateCreatePaymentRequest()

	payment, err := paymentService.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("user with id %s in not exists", request.UserId)), err)
	assert.Equal(t, model.PaymentResponse{}, payment)
	assert.Len(t, paymentService.(*paymentServiceObject).payments, 0)
}

func Test_Add_WithHouseNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", userId).Return(true)
	houses.On("ExistsById", houseId).Return(false)

	request := generateCreatePaymentRequest()

	payment, err := paymentService.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s in not exists", request.HouseId)), err)
	assert.Equal(t, model.PaymentResponse{}, payment)
	assert.Len(t, paymentService.(*paymentServiceObject).payments, 0)
}

func Test_FindPaymentById(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", userId).Return(true)
	houses.On("ExistsById", houseId).Return(true)

	request := generateCreatePaymentRequest()

	payment, err := paymentService.Add(request)

	assert.Nil(t, err)

	actual, err := paymentService.FindById(payment.Id)

	assert.Nil(t, nil)
	assert.Equal(t, payment, actual)
}

func Test_FindPaymentById_WithNotExistingId(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()
	actual, err := paymentService.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("payment with id %s not found", id)), err)
	assert.Equal(t, model.PaymentResponse{}, actual)
}

func Test_FindPaymentByHouseId(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", userId).Return(true)
	houses.On("ExistsById", houseId).Return(true)

	request := generateCreatePaymentRequest()

	payment, err := paymentService.Add(request)

	assert.Nil(t, err)

	payments := paymentService.FindByHouseId(payment.HouseId)

	assert.Equal(t, []model.PaymentResponse{payment}, payments)
}

func Test_FindPaymentByHouseId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	payments := paymentService.FindByHouseId(uuid.New())

	assert.Equal(t, []model.PaymentResponse{}, payments)
}

func Test_FindPaymentByUserId(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", userId).Return(true)
	houses.On("ExistsById", houseId).Return(true)

	request := generateCreatePaymentRequest()

	payment, err := paymentService.Add(request)

	assert.Nil(t, err)

	payments := paymentService.FindByUserId(payment.UserId)

	assert.Equal(t, []model.PaymentResponse{payment}, payments)
}

func Test_FindPaymentByUserId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	payments := paymentService.FindByUserId(uuid.New())

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
