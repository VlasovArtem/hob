package handler

import (
	"common/rest"
	"net/http"
	"payment/scheduler/model"
	"payment/scheduler/service"
)

type paymentSchedulerHandlerObject struct {
	paymentSchedulerService service.PaymentSchedulerService
}

func NewPaymentSchedulerHandler(paymentSchedulerService service.PaymentSchedulerService) PaymentSchedulerHandler {
	return &paymentSchedulerHandlerObject{paymentSchedulerService}
}

type PaymentSchedulerHandler interface {
	Add() http.HandlerFunc
	Remove() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
	FindByUserId() http.HandlerFunc
}

func (p *paymentSchedulerHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		paymentRequest := model.CreatePaymentSchedulerRequest{}

		if err := rest.PerformRequest(&paymentRequest, writer, request); err != nil {
			return
		}

		payment, err := p.paymentSchedulerService.Add(paymentRequest)

		rest.PerformResponseWithCode(writer, payment, http.StatusCreated, err)
	}
}

func (p *paymentSchedulerHandlerObject) Remove() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			err = p.paymentSchedulerService.Remove(id)

			rest.PerformResponseWithCode(writer, nil, http.StatusNoContent, err)
		}
	}
}

func (p *paymentSchedulerHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			response, err := p.paymentSchedulerService.FindById(id)

			rest.PerformResponse(writer, response, err)
		}
	}
}

func (p *paymentSchedulerHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			rest.PerformResponse(writer, p.paymentSchedulerService.FindByHouseId(id), nil)
		}
	}
}

func (p *paymentSchedulerHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			rest.PerformResponse(writer, p.paymentSchedulerService.FindByUserId(id), nil)
		}
	}
}