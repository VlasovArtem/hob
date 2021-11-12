package handler

import (
	"common/rest"
	"income/model"
	"income/service"
	"net/http"
)

type incomeHandlerObject struct {
	incomeService service.IncomeService
}

func NewIncomeHandler(incomeService service.IncomeService) IncomeHandler {
	return &incomeHandlerObject{incomeService}
}

type IncomeHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
}

func (i *incomeHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var requestBody model.CreateIncomeRequest

		if err := rest.PerformRequest(&requestBody, writer, request); err != nil {
			return
		}

		meter, err := i.incomeService.AddIncome(requestBody)

		rest.PerformResponseWithCode(writer, meter, http.StatusCreated, err)
	}
}

func (i *incomeHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			meterResponse, err := i.incomeService.FindById(id)

			rest.PerformResponse(writer, meterResponse, err)
		}
	}
}

func (i *incomeHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			if meterResponse, err := i.incomeService.FindByHouseId(id); err == nil {
				rest.PerformResponse(writer, meterResponse, err)
			}
		}
	}
}
