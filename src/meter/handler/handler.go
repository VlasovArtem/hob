package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/VlasovArtem/hob/src/meter/service"
	"github.com/gorilla/mux"
	"net/http"
)

type MeterHandlerStr struct {
	meterService service.MeterService
}

func NewMeterHandler(meterService service.MeterService) MeterHandler {
	return &MeterHandlerStr{meterService}
}

func (m *MeterHandlerStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(service.MeterServiceStr{}),
	}
}

func (m *MeterHandlerStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewMeterHandler(dependency.FindRequiredDependency[service.MeterServiceStr, service.MeterService](factory))
}

func (m *MeterHandlerStr) Init(router *mux.Router) {
	meterRouter := router.PathPrefix("/api/v1/meters").Subrouter()

	meterRouter.Path("").HandlerFunc(m.Add()).Methods("POST")
	meterRouter.Path("/{id}").HandlerFunc(m.FindById()).Methods("GET")
	meterRouter.Path("/{id}").HandlerFunc(m.FindById()).Methods("PUT")
	meterRouter.Path("/{id}").HandlerFunc(m.FindById()).Methods("DELETE")
	meterRouter.Path("/payment/{id}").HandlerFunc(m.FindByPaymentId()).Methods("GET")
}

type MeterHandler interface {
	Add() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByPaymentId() http.HandlerFunc
}

func (m *MeterHandlerStr) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if body, err := rest.ReadRequestBody[model.CreateMeterRequest](request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Created(m.meterService.Add(body)).
				Perform()
		}
	}
}

func (m *MeterHandlerStr) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Ok(m.meterService.FindById(id)).
				Perform()
		}
	}
}

func (m *MeterHandlerStr) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if body, err := rest.ReadRequestBody[model.UpdateMeterRequest](request); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				rest.NewAPIResponse(writer).
					Error(m.meterService.Update(id, body)).
					Perform()
			}
		}
	}
}

func (m *MeterHandlerStr) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				NoContent(m.meterService.DeleteById(id)).
				Perform()
		}
	}
}

func (m *MeterHandlerStr) FindByPaymentId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Ok(m.meterService.FindByPaymentId(id)).
				Perform()
		}
	}
}
