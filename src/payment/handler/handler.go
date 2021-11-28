package handler

import (
	"encoding/json"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/payment/model"
	paymentService "github.com/VlasovArtem/hob/src/payment/service"
	"github.com/gorilla/mux"
	"net/http"
)

type PaymentHandlerObject struct {
	paymentService paymentService.PaymentService
}

func NewPaymentHandler(paymentService paymentService.PaymentService) PaymentHandler {
	return &PaymentHandlerObject{
		paymentService: paymentService,
	}
}

func (p *PaymentHandlerObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewPaymentHandler(
			factory.FindRequiredByObject(paymentService.PaymentServiceObject{}).(paymentService.PaymentService),
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
