package handler

import (
	"common/rest"
	"income/scheduler/model"
	"income/scheduler/service"
	"net/http"
)

type incomeSchedulerHandlerObject struct {
	incomeSchedulerService service.IncomeSchedulerService
}

func NewIncomeSchedulerHandler(incomeSchedulerService service.IncomeSchedulerService) IncomeSchedulerHandler {
	return &incomeSchedulerHandlerObject{incomeSchedulerService}
}

type IncomeSchedulerHandler interface {
	Add() http.HandlerFunc
	Remove() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
}

func (i *incomeSchedulerHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		incomeRequest := model.CreateIncomeSchedulerRequest{}

		if err := rest.PerformRequest(&incomeRequest, writer, request); err != nil {
			return
		}

		income, err := i.incomeSchedulerService.Add(incomeRequest)

		rest.PerformResponseWithCode(writer, income, http.StatusCreated, err)
	}
}

func (i *incomeSchedulerHandlerObject) Remove() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			err = i.incomeSchedulerService.Remove(id)

			rest.PerformResponseWithCode(writer, nil, http.StatusNoContent, err)
		}
	}
}

func (i *incomeSchedulerHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			response, err := i.incomeSchedulerService.FindById(id)

			rest.PerformResponse(writer, response, err)
		}
	}
}

func (i *incomeSchedulerHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			response, err := i.incomeSchedulerService.FindByHouseId(id)

			rest.PerformResponse(writer, response, err)
		}
	}
}
