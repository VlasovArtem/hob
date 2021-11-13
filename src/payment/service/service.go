package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	hs "house/service"
	"payment/model"
	us "user/service"
)

type paymentServiceObject struct {
	userService   us.UserService
	houseService  hs.HouseService
	payments      map[uuid.UUID]model.Payment
	housePayments map[uuid.UUID][]model.Payment
	userPayments  map[uuid.UUID][]model.Payment
}

func NewPaymentService(userService us.UserService, houseService hs.HouseService) PaymentService {
	return &paymentServiceObject{
		userService:   userService,
		houseService:  houseService,
		payments:      make(map[uuid.UUID]model.Payment),
		housePayments: make(map[uuid.UUID][]model.Payment),
		userPayments:  make(map[uuid.UUID][]model.Payment),
	}
}

type PaymentService interface {
	Add(request model.CreatePaymentRequest) (model.PaymentResponse, error)
	FindById(id uuid.UUID) (model.PaymentResponse, error)
	FindByHouseId(houseId uuid.UUID) []model.PaymentResponse
	FindByUserId(userId uuid.UUID) []model.PaymentResponse
	ExistsById(id uuid.UUID) bool
}

func (p *paymentServiceObject) Add(request model.CreatePaymentRequest) (response model.PaymentResponse, err error) {
	if !p.userService.ExistsById(request.UserId) {
		return response, errors.New(fmt.Sprintf("user with id %s in not exists", request.UserId))
	}
	if !p.houseService.ExistsById(request.HouseId) {
		return response, errors.New(fmt.Sprintf("house with id %s in not exists", request.HouseId))
	}

	entity := request.ToEntity()

	p.payments[entity.Id] = entity
	p.housePayments[entity.HouseId] = append(p.housePayments[entity.HouseId], entity)
	p.userPayments[entity.UserId] = append(p.userPayments[entity.UserId], entity)

	return entity.ToResponse(), nil
}

func (p *paymentServiceObject) FindById(id uuid.UUID) (model.PaymentResponse, error) {
	if payment, ok := p.payments[id]; ok {
		return payment.ToResponse(), nil
	}
	return model.PaymentResponse{}, errors.New(fmt.Sprintf("payment with id %s not found", id))
}

func (p *paymentServiceObject) FindByHouseId(houseId uuid.UUID) []model.PaymentResponse {
	return convert(p.housePayments[houseId])
}

func (p *paymentServiceObject) FindByUserId(userId uuid.UUID) []model.PaymentResponse {
	return convert(p.userPayments[userId])
}

func (p *paymentServiceObject) ExistsById(id uuid.UUID) bool {
	_, ok := p.payments[id]

	return ok
}

func convert(payments []model.Payment) []model.PaymentResponse {
	if len(payments) == 0 {
		return make([]model.PaymentResponse, 0)
	}

	var paymentsResponse []model.PaymentResponse

	for _, payment := range payments {
		paymentsResponse = append(paymentsResponse, payment.ToResponse())
	}

	return paymentsResponse
}
