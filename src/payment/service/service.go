package service

import (
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/database"
	"github.com/VlasovArtem/hob/src/common/dependency"
	hs "github.com/VlasovArtem/hob/src/house/service"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/payment/repository"
	us "github.com/VlasovArtem/hob/src/user/service"
	"github.com/google/uuid"
)

type PaymentServiceObject struct {
	userService       us.UserService
	houseService      hs.HouseService
	paymentRepository repository.PaymentRepository
}

func NewPaymentService(userService us.UserService, houseService hs.HouseService, paymentRepository repository.PaymentRepository) PaymentService {
	return &PaymentServiceObject{
		userService:       userService,
		houseService:      houseService,
		paymentRepository: paymentRepository,
	}
}

func (p *PaymentServiceObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewPaymentService(
			factory.FindRequiredByObject(us.UserServiceObject{}).(us.UserService),
			factory.FindRequiredByObject(hs.HouseServiceObject{}).(hs.HouseService),
			factory.FindRequiredByObject(repository.PaymentRepositoryObject{}).(repository.PaymentRepository),
		),
	)
}

type PaymentService interface {
	Add(request model.CreatePaymentRequest) (model.PaymentDto, error)
	FindById(id uuid.UUID) (model.PaymentDto, error)
	FindByHouseId(houseId uuid.UUID) []model.PaymentDto
	FindByUserId(userId uuid.UUID) []model.PaymentDto
	ExistsById(id uuid.UUID) bool
}

func (p *PaymentServiceObject) Add(request model.CreatePaymentRequest) (response model.PaymentDto, err error) {
	if !p.userService.ExistsById(request.UserId) {
		return response, errors.New(fmt.Sprintf("user with id %s in not exists", request.UserId))
	}
	if !p.houseService.ExistsById(request.HouseId) {
		return response, errors.New(fmt.Sprintf("house with id %s in not exists", request.HouseId))
	} else if err != nil {
		return response, err
	}

	payment, err := p.paymentRepository.Create(request.ToEntity())

	return payment.ToDto(), err
}

func (p *PaymentServiceObject) FindById(id uuid.UUID) (model.PaymentDto, error) {
	if payment, err := p.paymentRepository.FindById(id); err != nil {
		return model.PaymentDto{}, database.HandlerFindError(err, fmt.Sprintf("payment with id %s not found", id))
	} else {
		return payment.ToDto(), nil
	}
}

func (p *PaymentServiceObject) FindByHouseId(houseId uuid.UUID) []model.PaymentDto {
	return convert(p.paymentRepository.FindByHouseId(houseId))
}

func (p *PaymentServiceObject) FindByUserId(userId uuid.UUID) []model.PaymentDto {
	return convert(p.paymentRepository.FindByUserId(userId))
}

func (p *PaymentServiceObject) ExistsById(id uuid.UUID) bool {
	return p.paymentRepository.ExistsById(id)
}

func convert(payments []model.Payment) []model.PaymentDto {
	if len(payments) == 0 {
		return make([]model.PaymentDto, 0)
	}

	var paymentsResponse []model.PaymentDto

	for _, payment := range payments {
		paymentsResponse = append(paymentsResponse, payment.ToDto())
	}

	return paymentsResponse
}
