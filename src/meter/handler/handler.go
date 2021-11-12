package handler

import (
	"common/rest"
	"meter/model"
	"meter/service"
	"net/http"
)

type meterHandlerObject struct {
	meterService service.MeterService
}

func NewMeterHandler(meterService service.MeterService) MeterHandler {
	return &meterHandlerObject{meterService}
}

type MeterHandler interface {
	AddMeter() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByPaymentId() http.HandlerFunc
}

func (m *meterHandlerObject) AddMeter() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var requestBody model.CreateMeterRequest

		if err := rest.PerformRequest(&requestBody, writer, request); err != nil {
			return
		}

		meter, err := m.meterService.AddMeter(requestBody)

		rest.PerformResponseWithCode(writer, meter, http.StatusCreated, err)
	}
}

func (m *meterHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			meterResponse, err := m.meterService.FindById(id)

			rest.PerformResponse(writer, meterResponse, err)
		}
	}
}

func (m *meterHandlerObject) FindByPaymentId() http.HandlerFunc {
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
