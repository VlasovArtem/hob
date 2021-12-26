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

func (p *PaymentHandlerObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewPaymentHandler(factory.FindRequiredByObject(paymentService.PaymentServiceObject{}).(paymentService.PaymentService))
}

func (p *PaymentHandlerObject) Init(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1/payment").Subrouter()

	subrouter.Path("").HandlerFunc(p.Add()).Methods("POST")
	subrouter.Path("/{id}").HandlerFunc(p.FindById()).Methods("GET")
	subrouter.Path("/{id}").HandlerFunc(p.Delete()).Methods("DELETE")
	subrouter.Path("/{id}").HandlerFunc(p.Update()).Methods("PUT")
	subrouter.Path("/house/{id}").HandlerFunc(p.FindByHouseId()).Methods("GET")
	subrouter.Path("/user/{id}").HandlerFunc(p.FindByUserId()).Methods("GET")
	subrouter.Path("/provider/{id}").HandlerFunc(p.FindByProviderId()).Methods("GET")
}

type PaymentHandler interface {
	Add() http.HandlerFunc
	Delete() http.HandlerFunc
	Update() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
	FindByUserId() http.HandlerFunc
	FindByProviderId() http.HandlerFunc
}

func (p *PaymentHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		paymentRequest := model.CreatePaymentRequest{}

		if err := rest.ReadRequestBody(&paymentRequest, writer, request); err != nil {
			return
		}

		if payment, err := p.paymentService.Add(paymentRequest); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			writer.WriteHeader(http.StatusCreated)
			json.NewEncoder(writer).Encode(payment)
		}
	}
}

func (p *PaymentHandlerObject) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.PerformResponseWithCode(writer, nil, http.StatusNoContent, p.paymentService.DeleteById(id))
		}
	}
}

func (p *PaymentHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		paymentRequest := model.UpdatePaymentRequest{}

		if err := rest.ReadRequestBody(&paymentRequest, writer, request); err != nil {
			return
		}

		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.PerformResponseWithCode(writer, nil, http.StatusNoContent, p.paymentService.Update(id, paymentRequest))
		}
	}
}

func (p *PaymentHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			response, err := p.paymentService.FindById(id)

			rest.PerformResponseWithBody(writer, response, err)
		}
	}
}

func (p *PaymentHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.PerformResponseWithBody(writer, p.paymentService.FindByHouseId(id), nil)
		}
	}
}

func (p *PaymentHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.PerformResponseWithBody(writer, p.paymentService.FindByUserId(id), nil)
		}
	}
}

func (p *PaymentHandlerObject) FindByProviderId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.PerformResponseWithBody(writer, p.paymentService.FindByProviderId(id), nil)
		}
	}
}
