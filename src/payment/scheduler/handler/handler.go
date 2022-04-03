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

func (p *PaymentSchedulerHandlerObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewPaymentSchedulerHandler(dependency.FindRequiredDependency[service.PaymentSchedulerServiceObject, service.PaymentSchedulerService](factory))
}

func (p *PaymentSchedulerHandlerObject) Init(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1/payment/scheduler").Subrouter()

	subrouter.Path("").HandlerFunc(p.Add()).Methods("POST")
	subrouter.Path("/{id}").HandlerFunc(p.FindById()).Methods("GET")
	subrouter.Path("/{id}").HandlerFunc(p.Remove()).Methods("DELETE")
	subrouter.Path("/{id}").HandlerFunc(p.Update()).Methods("PUT")
	subrouter.Path("/house/{id}").HandlerFunc(p.FindByHouseId()).Methods("GET")
	subrouter.Path("/user/{id}").HandlerFunc(p.FindByUserId()).Methods("GET")
	subrouter.Path("/provider/{id}").HandlerFunc(p.FindByUserId()).Methods("GET")
}

type PaymentSchedulerHandler interface {
	Add() http.HandlerFunc
	Remove() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
	FindByUserId() http.HandlerFunc
	FindByProviderId() http.HandlerFunc
	Update() http.HandlerFunc
}

func (p *PaymentSchedulerHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if body, err := rest.ReadRequestBody[model.CreatePaymentSchedulerRequest](request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Created(p.paymentSchedulerService.Add(body)).
				Perform()
		}
	}
}

func (p *PaymentSchedulerHandlerObject) Remove() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				NoContent(p.paymentSchedulerService.Remove(id)).
				Perform()
		}
	}
}

func (p *PaymentSchedulerHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if body, err := rest.ReadRequestBody[model.UpdatePaymentSchedulerRequest](request); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				rest.NewAPIResponse(writer).
					Error(p.paymentSchedulerService.Update(id, body)).
					Perform()
			}
		}
	}
}

func (p *PaymentSchedulerHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Ok(p.paymentSchedulerService.FindById(id)).
				Perform()
		}
	}
}

func (p *PaymentSchedulerHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Body(p.paymentSchedulerService.FindByHouseId(id)).
				Perform()
		}
	}
}

func (p *PaymentSchedulerHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Body(p.paymentSchedulerService.FindByUserId(id)).
				Perform()
		}
	}
}

func (p *PaymentSchedulerHandlerObject) FindByProviderId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Body(p.paymentSchedulerService.FindByProviderId(id)).
				Perform()
		}
	}
}
