package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/meter/model"
	"github.com/VlasovArtem/hob/src/meter/service"
	"github.com/gorilla/mux"
	"net/http"
)

type MeterHandlerObject struct {
	meterService service.MeterService
}

func NewMeterHandler(meterService service.MeterService) MeterHandler {
	return &MeterHandlerObject{meterService}
}

func (m *MeterHandlerObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewMeterHandler(factory.FindRequiredByType(service.MeterServiceType).(service.MeterService))
}

func (m *MeterHandlerObject) Init(router *mux.Router) {
	meterRouter := router.PathPrefix("/api/v1/meter").Subrouter()

	meterRouter.Path("").HandlerFunc(m.Add()).Methods("POST")
	meterRouter.Path("/{id}").HandlerFunc(m.FindById()).Methods("GET")
	meterRouter.Path("/{id}").HandlerFunc(m.FindById()).Methods("PUT")
	meterRouter.Path("/{id}").HandlerFunc(m.FindById()).Methods("DELETE")
	meterRouter.Path("/payment/{id}").HandlerFunc(m.FindByPaymentId()).Methods("GET")
	meterRouter.Path("/house/{id}").HandlerFunc(m.FindByHouseId()).Methods("GET")
}

type MeterHandler interface {
	Add() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByPaymentId() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
}

func (m *MeterHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var requestBody model.CreateMeterRequest

		if err := rest.ReadRequestBody(&requestBody, writer, request); err != nil {
			return
		}

		meter, err := m.meterService.Add(requestBody)

		rest.PerformResponseWithCode(writer, meter, http.StatusCreated, err)
	}
}

func (m *MeterHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			meterResponse, err := m.meterService.FindById(id)

			rest.PerformResponseWithBody(writer, meterResponse, err)
		}
	}
}

func (m *MeterHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			var requestBody model.UpdateMeterRequest

			if err := rest.ReadRequestBody(&requestBody, writer, request); err != nil {
				return
			}

			rest.PerformResponseWithBody(writer, nil, m.meterService.Update(id, requestBody))
		}
	}
}

func (m *MeterHandlerObject) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {

			rest.PerformResponseWithCode(writer, nil, http.StatusNoContent, m.meterService.DeleteById(id))
		}
	}
}

func (m *MeterHandlerObject) FindByPaymentId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if meterResponse, err := m.meterService.FindByPaymentId(id); err == nil {
				rest.PerformResponseWithBody(writer, meterResponse, err)
			}
		}
	}
}

func (m *MeterHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.PerformResponseWithBody(writer, m.meterService.FindByHouseId(id), err)
		}
	}
}
