package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/payment/scheduler/model"
	"github.com/VlasovArtem/hob/src/payment/scheduler/service"
	"github.com/gorilla/mux"
	"net/http"
)

type PaymentSchedulerHandlerObject struct {
	paymentSchedulerService service.PaymentSchedulerService
}

func NewPaymentSchedulerHandler(paymentSchedulerService service.PaymentSchedulerService) PaymentSchedulerHandler {
	return &PaymentSchedulerHandlerObject{paymentSchedulerService}
}

func (p *PaymentSchedulerHandlerObject) Initialize(factory dependency.DependenciesFactory) interface{} {
	return factory.Add(
		NewPaymentSchedulerHandler(factory.FindRequiredByObject(service.PaymentSchedulerServiceObject{}).(service.PaymentSchedulerService)),
	)
}

func (p *PaymentSchedulerHandlerObject) Init(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1/payment/scheduler").Subrouter()

	subrouter.Path("").HandlerFunc(p.Add()).Methods("POST")
	subrouter.Path("/{id}").HandlerFunc(p.FindById()).Methods("GET")
	subrouter.Path("/{id}").HandlerFunc(p.Remove()).Methods("DELETE")
	subrouter.Path("/house/{id}").HandlerFunc(p.FindByHouseId()).Methods("GET")
	subrouter.Path("/user/{id}").HandlerFunc(p.FindByUserId()).Methods("GET")
}

type PaymentSchedulerHandler interface {
	Add() http.HandlerFunc
	Remove() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
	FindByUserId() http.HandlerFunc
}

func (p *PaymentSchedulerHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		paymentRequest := model.CreatePaymentSchedulerRequest{}

		if err := rest.PerformRequest(&paymentRequest, writer, request); err != nil {
			return
		}

		payment, err := p.paymentSchedulerService.Add(paymentRequest)

		rest.PerformResponseWithCode(writer, payment, http.StatusCreated, err)
	}
}

func (p *PaymentSchedulerHandlerObject) Remove() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			err = p.paymentSchedulerService.Remove(id)

			rest.PerformResponseWithCode(writer, nil, http.StatusNoContent, err)
		}
	}
}

func (p *PaymentSchedulerHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			response, err := p.paymentSchedulerService.FindById(id)

			rest.PerformResponse(writer, response, err)
		}
	}
}

func (p *PaymentSchedulerHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			rest.PerformResponse(writer, p.paymentSchedulerService.FindByHouseId(id), nil)
		}
	}
}

func (p *PaymentSchedulerHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			rest.PerformResponse(writer, p.paymentSchedulerService.FindByUserId(id), nil)
		}
	}
}
