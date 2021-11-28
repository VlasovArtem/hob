package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	houseMocks "house/mocks"
	"payment/mocks"
	"payment/model"
	"testing"
	userMocks "user/mocks"
)

var (
	users             *userMocks.UserService
	houses            *houseMocks.HouseService
	paymentRepository *mocks.PaymentRepository
)

func serviceGenerator() PaymentService {
	users = new(userMocks.UserService)
	houses = new(houseMocks.HouseService)
	paymentRepository = new(mocks.PaymentRepository)

	return NewPaymentService(users, houses, paymentRepository)
}

func Test_Add(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", mocks.UserId).Return(true)
	houses.On("ExistsById", mocks.HouseId).Return(true)
	paymentRepository.On("Create", mock.Anything).Return(
		func(payment model.Payment) model.Payment { return payment },
		nil,
	)

	request := mocks.GenerateCreatePaymentRequest()

	payment, err := paymentService.Add(request)

	expectedEntity := request.ToEntity()
	expectedEntity.Id = payment.Id
	expectedResponse := expectedEntity.ToDto()

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, payment)
}

func Test_Add_WithUserNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", mocks.UserId).Return(false)

	request := mocks.GenerateCreatePaymentRequest()

	payment, err := paymentService.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("user with id %s in not exists", request.UserId)), err)
	assert.Equal(t, model.PaymentDto{}, payment)
}

func Test_Add_WithHouseNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	users.On("ExistsById", mocks.UserId).Return(true)
	houses.On("ExistsById", mocks.HouseId).Return(false)

	request := mocks.GenerateCreatePaymentRequest()

	payment, err := paymentService.Add(request)

	assert.Equal(t, errors.New(fmt.Sprintf("house with id %s in not exists", request.HouseId)), err)
	assert.Equal(t, model.PaymentDto{}, payment)
}

func Test_FindById(t *testing.T) {
	paymentService := serviceGenerator()

	payment := mocks.GeneratePayment(mocks.HouseId, mocks.UserId)

	paymentRepository.On("FindById", payment.Id).Return(payment, nil)

	actual, err := paymentService.FindById(payment.Id)

	assert.Nil(t, err)
	assert.Equal(t, payment.ToDto(), actual)
}

func Test_FindById_WithNotExistingId(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	paymentRepository.On("FindById", id).Return(model.Payment{}, gorm.ErrRecordNotFound)

	actual, err := paymentService.FindById(id)

	assert.Equal(t, errors.New(fmt.Sprintf("payment with id %s not found", id)), err)
	assert.Equal(t, model.PaymentDto{}, actual)
}

func Test_FindById_WithError(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()
	err2 := errors.New("error")

	paymentRepository.On("FindById", id).Return(model.Payment{}, err2)

	actual, err := paymentService.FindById(id)

	assert.Equal(t, err2, err)
	assert.Equal(t, model.PaymentDto{}, actual)
}

func Test_FindByHouseId(t *testing.T) {
	paymentService := serviceGenerator()

	payment := mocks.GeneratePayment(mocks.HouseId, mocks.UserId)

	houseId := uuid.New()

	paymentRepository.On("FindByHouseId", houseId).Return([]model.Payment{payment})

	payments := paymentService.FindByHouseId(houseId)

	assert.Equal(t, []model.PaymentDto{payment.ToDto()}, payments)
}

func Test_FindByHouseId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	houseId := uuid.New()

	paymentRepository.On("FindByHouseId", houseId).Return([]model.Payment{})

	payments := paymentService.FindByHouseId(houseId)

	assert.Equal(t, []model.PaymentDto{}, payments)
}

func Test_FindByUserId(t *testing.T) {
	paymentService := serviceGenerator()

	payment := mocks.GeneratePayment(mocks.HouseId, mocks.UserId)

	userId := uuid.New()

	paymentRepository.On("FindByUserId", userId).Return([]model.Payment{payment})

	payments := paymentService.FindByUserId(userId)

	assert.Equal(t, []model.PaymentDto{payment.ToDto()}, payments)
}

func Test_FindByUserId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	userId := uuid.New()

	paymentRepository.On("FindByUserId", userId).Return([]model.Payment{})

	payments := paymentService.FindByUserId(userId)

	assert.Equal(t, []model.PaymentDto{}, payments)
}

func Test_ExistsById(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	paymentRepository.On("ExistsById", id).Return(true)

	assert.True(t, paymentService.ExistsById(id))
}

func Test_ExistsById_WithNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	paymentRepository.On("ExistsById", id).Return(false)

	assert.False(t, paymentService.ExistsById(id))
}
