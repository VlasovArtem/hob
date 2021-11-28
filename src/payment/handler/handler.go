package handler

import (
	"common/dependency"
	"common/rest"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"payment/model"
	ps "payment/service"
)

type PaymentHandlerObject struct {
	paymentService ps.PaymentService
}

func NewPaymentHandler(paymentService ps.PaymentService) PaymentHandler {
	return &PaymentHandlerObject{
		paymentService: paymentService,
	}
}

func (p *PaymentHandlerObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewPaymentHandler(
			factory.FindRequiredByObject(ps.PaymentServiceObject{}).(ps.PaymentService),
		),
	)
}

func (p *PaymentHandlerObject) Init(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1/payment").Subrouter()

	subrouter.Path("").HandlerFunc(p.Add()).Methods("POST")
	subrouter.Path("/{id}").HandlerFunc(p.FindById()).Methods("GET")
	subrouter.Path("/house/{id}").HandlerFunc(p.FindByHouseId()).Methods("GET")
	subrouter.Path("/user/{id}").HandlerFunc(p.FindByUserId()).Methods("GET")
}

type PaymentHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
	FindByUserId() http.HandlerFunc
}

func (p *PaymentHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		paymentRequest := model.CreatePaymentRequest{}

		if err := rest.PerformRequest(&paymentRequest, writer, request); err != nil {
			return
		}

		if payment, err := p.paymentService.Add(paymentRequest); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			writer.WriteHeader(http.StatusCreated)
			json.NewEncoder(writer).Encode(payment)
		}
	}
}

func (p *PaymentHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			response, err := p.paymentService.FindById(id)

			rest.PerformResponse(writer, response, err)
		}
	}
}

func (p *PaymentHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			rest.PerformResponse(writer, p.paymentService.FindByHouseId(id), nil)
		}
	}
}

func (p *PaymentHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			rest.PerformResponse(writer, p.paymentService.FindByUserId(id), nil)
		}
	}
}
