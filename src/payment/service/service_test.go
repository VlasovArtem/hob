package service

import (
	"errors"
	"fmt"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	"github.com/VlasovArtem/hob/src/payment/mocks"
	"github.com/VlasovArtem/hob/src/payment/model"
	providerMocks "github.com/VlasovArtem/hob/src/provider/mocks"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
	"time"
)

var (
	userService       *userMocks.UserService
	houseService      *houseMocks.HouseService
	providerService   *providerMocks.ProviderService
	paymentRepository *mocks.PaymentRepository
)

func serviceGenerator() PaymentService {
	userService = new(userMocks.UserService)
	houseService = new(houseMocks.HouseService)
	providerService = new(providerMocks.ProviderService)
	paymentRepository = new(mocks.PaymentRepository)

	return NewPaymentService(userService, houseService, providerService, paymentRepository)
}

func Test_Add(t *testing.T) {
	paymentService := serviceGenerator()

	userService.On("ExistsById", mocks.UserId).Return(true)
	houseService.On("ExistsById", mocks.HouseId).Return(true)
	providerService.On("ExistsById", mocks.ProviderId).Return(true)
	paymentRepository.On("Create", mock.Anything).Return(
		func(payment model.Payment) model.Payment { return payment },
		nil,
	)

	request := mocks.GenerateCreatePaymentRequest()

	payment, err := paymentService.Add(request)

	expectedEntity := request.CreateToEntity()
	expectedEntity.Id = payment.Id
	expectedResponse := expectedEntity.ToDto()

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, payment)
}

func Test_Add_WithUserNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	userService.On("ExistsById", mocks.UserId).Return(false)

	request := mocks.GenerateCreatePaymentRequest()

	payment, err := paymentService.Add(request)

	assert.Equal(t, fmt.Errorf("user with id %s in not exists", request.UserId), err)
	assert.Equal(t, model.PaymentDto{}, payment)
}

func Test_Add_WithHouseNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	userService.On("ExistsById", mocks.UserId).Return(true)
	houseService.On("ExistsById", mocks.HouseId).Return(false)

	request := mocks.GenerateCreatePaymentRequest()

	payment, err := paymentService.Add(request)

	assert.Equal(t, fmt.Errorf("house with id %s in not exists", request.HouseId), err)
	assert.Equal(t, model.PaymentDto{}, payment)
}

func Test_Add_WithProviderNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	userService.On("ExistsById", mocks.UserId).Return(true)
	houseService.On("ExistsById", mocks.HouseId).Return(true)
	providerService.On("ExistsById", mocks.ProviderId).Return(false)

	request := mocks.GenerateCreatePaymentRequest()

	payment, err := paymentService.Add(request)

	assert.Equal(t, fmt.Errorf("provider with id %s in not exists", request.ProviderId), err)
	assert.Equal(t, model.PaymentDto{}, payment)

	paymentRepository.AssertNotCalled(t, "Create", mock.Anything)
}

func Test_FindById(t *testing.T) {
	paymentService := serviceGenerator()

	payment := mocks.GeneratePayment(mocks.HouseId, mocks.UserId, mocks.ProviderId)

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

	assert.Equal(t, fmt.Errorf("payment with id %s not found", id), err)
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

	payment := mocks.GeneratePayment(mocks.HouseId, mocks.UserId, mocks.ProviderId)

	houseId := uuid.New()

	dto := payment.ToDto()
	paymentRepository.On("FindByHouseId", houseId).Return([]model.PaymentDto{dto})

	payments := paymentService.FindByHouseId(houseId)

	assert.Equal(t, []model.PaymentDto{dto}, payments)
}

func Test_FindByHouseId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	houseId := uuid.New()

	paymentRepository.On("FindByHouseId", houseId).Return([]model.PaymentDto{})

	payments := paymentService.FindByHouseId(houseId)

	assert.Equal(t, []model.PaymentDto{}, payments)
}

func Test_FindByUserId(t *testing.T) {
	paymentService := serviceGenerator()

	payment := mocks.GeneratePayment(mocks.HouseId, mocks.UserId, mocks.ProviderId)

	userId := uuid.New()

	dto := payment.ToDto()
	paymentRepository.On("FindByUserId", userId).Return([]model.PaymentDto{dto})

	payments := paymentService.FindByUserId(userId)

	assert.Equal(t, []model.PaymentDto{dto}, payments)
}

func Test_FindByUserId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	userId := uuid.New()

	paymentRepository.On("FindByUserId", userId).Return([]model.PaymentDto{})

	payments := paymentService.FindByUserId(userId)

	assert.Equal(t, []model.PaymentDto{}, payments)
}

func Test_FindByProviderId(t *testing.T) {
	paymentService := serviceGenerator()

	payment := mocks.GeneratePayment(mocks.HouseId, mocks.UserId, mocks.ProviderId)

	userId := uuid.New()

	dto := payment.ToDto()
	paymentRepository.On("FindByProviderId", userId).Return([]model.PaymentDto{dto})

	payments := paymentService.FindByProviderId(userId)

	assert.Equal(t, []model.PaymentDto{dto}, payments)
}

func Test_FindByProviderId_WithNotExistingRecords(t *testing.T) {
	paymentService := serviceGenerator()

	userId := uuid.New()

	paymentRepository.On("FindByProviderId", userId).Return([]model.PaymentDto{})

	payments := paymentService.FindByProviderId(userId)

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

func Test_DeleteById(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	paymentRepository.On("ExistsById", id).Return(true)
	paymentRepository.On("DeleteById", id).Return(nil)

	assert.Nil(t, paymentService.DeleteById(id))
}

func Test_DeleteById_WithNotExists(t *testing.T) {
	paymentService := serviceGenerator()

	id := uuid.New()

	paymentRepository.On("ExistsById", id).Return(false)

	assert.Equal(t, fmt.Errorf("payment with id %s not found", id), paymentService.DeleteById(id))

	paymentRepository.AssertNotCalled(t, "DeleteById", id)
}

func Test_Update(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdatePaymentRequest()
	id := uuid.New()

	paymentRepository.On("ExistsById", id).Return(true)
	providerService.On("ExistsById", request.ProviderId).Return(true)
	paymentRepository.On("Update", mock.Anything).Return(nil)

	assert.Nil(t, houseService.Update(id, request))

	paymentRepository.AssertCalled(t, "Update", model.Payment{
		Id:          id,
		Name:        request.Name,
		Description: request.Description,
		ProviderId:  request.ProviderId,
		Date:        request.Date,
		Sum:         request.Sum,
	})
}

func Test_Update_WithErrorFromDatabase(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdatePaymentRequest()
	id := uuid.New()

	paymentRepository.On("ExistsById", id).Return(true)
	providerService.On("ExistsById", request.ProviderId).Return(true)
	paymentRepository.On("Update", mock.Anything).Return(errors.New("test"))

	err := houseService.Update(id, request)
	assert.Equal(t, errors.New("test"), err)
}

func Test_Update_WithNotExists(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdatePaymentRequest()
	id := uuid.New()

	paymentRepository.On("ExistsById", id).Return(false)

	err := houseService.Update(id, request)
	assert.Equal(t, fmt.Errorf("payment with id %s not found", id), err)

	paymentRepository.AssertNotCalled(t, "Update", mock.Anything)
}

func Test_Update_WithDateAfterCurrentDate(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdatePaymentRequest()
	request.Date = time.Now().Add(time.Hour)
	id := uuid.New()

	paymentRepository.On("ExistsById", id).Return(true)
	providerService.On("ExistsById", request.ProviderId).Return(true)

	err := houseService.Update(id, request)
	assert.Equal(t, errors.New("date should not be after current date"), err)

	paymentRepository.AssertNotCalled(t, "Update", mock.Anything)
}

func Test_Update_WithProviderNotExists(t *testing.T) {
	houseService := serviceGenerator()

	request := mocks.GenerateUpdatePaymentRequest()
	id := uuid.New()

	paymentRepository.On("ExistsById", id).Return(true)
	providerService.On("ExistsById", request.ProviderId).Return(false)

	err := houseService.Update(id, request)
	assert.Equal(t, fmt.Errorf("provider with id %s not found", request.ProviderId), err)

	paymentRepository.AssertNotCalled(t, "Update", mock.Anything)
}
