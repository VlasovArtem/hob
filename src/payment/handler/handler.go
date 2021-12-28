package handler

import (
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

func (p *PaymentHandlerObject) Initialize(factory dependency.DependenciesProvider) any {
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
		if body, err := rest.ReadRequestBody[model.CreatePaymentRequest](request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Created(p.paymentService.Add(body)).
				Perform()
		}
	}
}

func (p *PaymentHandlerObject) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				NoContent(p.paymentService.DeleteById(id)).
				Perform()
		}
	}
}

func (p *PaymentHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if body, err := rest.ReadRequestBody[model.UpdatePaymentRequest](request); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				rest.NewAPIResponse(writer).
					Error(p.paymentService.Update(id, body)).
					Perform()
			}
		}
	}
}

func (p *PaymentHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Ok(p.paymentService.FindById(id)).
				Perform()
		}
	}
}

func (p *PaymentHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Body(p.paymentService.FindByHouseId(id)).
				Perform()
		}
	}
}

func (p *PaymentHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Body(p.paymentService.FindByUserId(id)).
				Perform()
		}
	}
}

func (p *PaymentHandlerObject) FindByProviderId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Body(p.paymentService.FindByProviderId(id)).
				Perform()
		}
	}
}
