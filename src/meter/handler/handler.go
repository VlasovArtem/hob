package handler

import (
	"common/dependency"
	"common/rest"
	"github.com/gorilla/mux"
	"meter/model"
	"meter/service"
	"net/http"
)

type MeterHandlerObject struct {
	meterService service.MeterService
}

func NewMeterHandler(meterService service.MeterService) MeterHandler {
	return &MeterHandlerObject{meterService}
}

func (m *MeterHandlerObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(NewMeterHandler(factory.FindRequiredByObject(service.MeterServiceObject{}).(service.MeterService)))
}

func (m *MeterHandlerObject) Init(router *mux.Router) {
	meterRouter := router.PathPrefix("/api/v1/meter").Subrouter()

	meterRouter.Path("").HandlerFunc(m.Add()).Methods("POST")
	meterRouter.Path("/{id}").HandlerFunc(m.FindById()).Methods("GET")
	meterRouter.Path("/payment/{id}").HandlerFunc(m.FindByPaymentId()).Methods("GET")
	meterRouter.Path("/house/{id}").HandlerFunc(m.FindByHouseId()).Methods("GET")
}

type MeterHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByPaymentId() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
}

func (m *MeterHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var requestBody model.CreateMeterRequest

		if err := rest.PerformRequest(&requestBody, writer, request); err != nil {
			return
		}

		meter, err := m.meterService.Add(requestBody)

		rest.PerformResponseWithCode(writer, meter, http.StatusCreated, err)
	}
}

func (m *MeterHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			meterResponse, err := m.meterService.FindById(id)

			rest.PerformResponse(writer, meterResponse, err)
		}
	}
}

func (m *MeterHandlerObject) FindByPaymentId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			if meterResponse, err := m.meterService.FindByPaymentId(id); err == nil {
				rest.PerformResponse(writer, meterResponse, err)
			}
		}
	}
}

func (m *MeterHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			rest.PerformResponse(writer, m.meterService.FindByHouseId(id), err)
		}
	}
}
