package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	hs "house/service"
	"log"
	pm "payment/model"
	"payment/scheduler/model"
	ps "payment/service"
	"scheduler"
	"time"
	us "user/service"
)

type paymentSchedulerServiceObject struct {
	userService      us.UserService
	houseService     hs.HouseService
	paymentService   ps.PaymentService
	serviceScheduler scheduler.ServiceScheduler
	payments         map[uuid.UUID]model.PaymentScheduler
	housePayments    map[uuid.UUID][]model.PaymentScheduler
	userPayments     map[uuid.UUID][]model.PaymentScheduler
}

func NewPaymentSchedulerService(
	userService us.UserService,
	houseService hs.HouseService,
	paymentService ps.PaymentService,
	serviceScheduler scheduler.ServiceScheduler,
) PaymentSchedulerService {
	return &paymentSchedulerServiceObject{
		userService:      userService,
		houseService:     houseService,
		paymentService:   paymentService,
		serviceScheduler: serviceScheduler,
		payments:         make(map[uuid.UUID]model.PaymentScheduler),
		housePayments:    make(map[uuid.UUID][]model.PaymentScheduler),
		userPayments:     make(map[uuid.UUID][]model.PaymentScheduler),
	}
}

type PaymentSchedulerService interface {
	Add(request model.CreatePaymentSchedulerRequest) (model.PaymentSchedulerResponse, error)
	Remove(id uuid.UUID) error
	FindById(id uuid.UUID) (model.PaymentSchedulerResponse, error)
	FindByHouseId(houseId uuid.UUID) []model.PaymentSchedulerResponse
	FindByUserId(userId uuid.UUID) []model.PaymentSchedulerResponse
}

func (p *paymentSchedulerServiceObject) Add(request model.CreatePaymentSchedulerRequest) (response model.PaymentSchedulerResponse, err error) {
	if request.Sum <= 0 {
		return response, errors.New("sum should not be zero of negative")
	}
	if !p.userService.ExistsById(request.UserId) {
		return response, errors.New(fmt.Sprintf("user with id %s in not exists", request.UserId))
	}
	if !p.houseService.ExistsById(request.HouseId) {
		return response, errors.New(fmt.Sprintf("house with id %s in not exists", request.HouseId))
	}
	if request.Spec == "" {
		return response, errors.New("scheduler configuration not provided")
	}

	entity := request.ToEntity()

	if _, err = p.serviceScheduler.Add(entity.Id, string(entity.Spec), p.schedulerFunc(entity.Payment)); err != nil {
		return response, err
	}

	p.payments[entity.Id] = entity
	p.housePayments[entity.HouseId] = append(p.housePayments[entity.HouseId], entity)
	p.userPayments[entity.UserId] = append(p.userPayments[entity.UserId], entity)

	return entity.ToResponse(), nil
}

func (p *paymentSchedulerServiceObject) Remove(id uuid.UUID) error {
	if paymentScheduler, ok := p.payments[id]; !ok {
		return errors.New(fmt.Sprintf("payment scheduler with id %s not found", id))
	} else {
		if err := p.serviceScheduler.Remove(id); err != nil {
			log.Println(err)
		}
		delete(p.payments, id)
		delete(p.userPayments, paymentScheduler.UserId)
		delete(p.housePayments, paymentScheduler.HouseId)
	}
	return nil
}

func (p *paymentSchedulerServiceObject) FindById(id uuid.UUID) (model.PaymentSchedulerResponse, error) {
	if payment, ok := p.payments[id]; ok {
		return payment.ToResponse(), nil
	}
	return model.PaymentSchedulerResponse{}, errors.New(fmt.Sprintf("payment scheduler with id %s not found", id))
}

func (p *paymentSchedulerServiceObject) FindByHouseId(houseId uuid.UUID) []model.PaymentSchedulerResponse {
	return convert(p.housePayments[houseId])
}

func (p *paymentSchedulerServiceObject) FindByUserId(userId uuid.UUID) []model.PaymentSchedulerResponse {
	return convert(p.userPayments[userId])
}

func (p *paymentSchedulerServiceObject) schedulerFunc(payment pm.Payment) func() {
	return func() {
		if _, err := p.paymentService.Add(
			pm.CreatePaymentRequest{
				Name:        payment.Name,
				Description: payment.Description,
				HouseId:     payment.HouseId,
				UserId:      payment.UserId,
				Date:        time.Now(),
				Sum:         payment.Sum,
			},
		); err != nil {
			log.Println(err)
		} else {
			log.Println(fmt.Sprintf("New payment added to the house %s and user %s via scheduler %s", payment.HouseId, payment.UserId, payment.Id))
		}
	}
}

func convert(payments []model.PaymentScheduler) []model.PaymentSchedulerResponse {
	if len(payments) == 0 {
		return make([]model.PaymentSchedulerResponse, 0)
	}

	var paymentsResponse []model.PaymentSchedulerResponse

	for _, payment := range payments {
		paymentsResponse = append(paymentsResponse, payment.ToResponse())
	}

	return paymentsResponse
}
