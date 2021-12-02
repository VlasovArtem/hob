package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/income/scheduler/model"
	"github.com/VlasovArtem/hob/src/income/scheduler/service"
	"github.com/gorilla/mux"
	"net/http"
)

type IncomeSchedulerHandlerObject struct {
	incomeSchedulerService service.IncomeSchedulerService
}

func NewIncomeSchedulerHandler(incomeSchedulerService service.IncomeSchedulerService) IncomeSchedulerHandler {
	return &IncomeSchedulerHandlerObject{incomeSchedulerService}
}

func (i *IncomeSchedulerHandlerObject) Initialize(factory dependency.DependenciesFactory) interface{} {
	return factory.Add(NewIncomeSchedulerHandler(factory.FindRequiredByObject(service.IncomeSchedulerServiceObject{}).(service.IncomeSchedulerService)))
}

func (i *IncomeSchedulerHandlerObject) Init(router *mux.Router) {
	incomeSchedulerRouter := router.PathPrefix("/api/v1/income/scheduler").Subrouter()

	incomeSchedulerRouter.Path("").HandlerFunc(i.Add()).Methods("POST")
	incomeSchedulerRouter.Path("/{id}").HandlerFunc(i.FindById()).Methods("GET")
	incomeSchedulerRouter.Path("/{id}").HandlerFunc(i.Remove()).Methods("DELETE")
	incomeSchedulerRouter.Path("/house/{id}").HandlerFunc(i.FindByHouseId()).Methods("GET")
}

type IncomeSchedulerHandler interface {
	Add() http.HandlerFunc
	Remove() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
}

func (i *IncomeSchedulerHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		incomeRequest := model.CreateIncomeSchedulerRequest{}

		if err := rest.PerformRequest(&incomeRequest, writer, request); err != nil {
			return
		}

		income, err := i.incomeSchedulerService.Add(incomeRequest)

		rest.PerformResponseWithCode(writer, income, http.StatusCreated, err)
	}
}

func (i *IncomeSchedulerHandlerObject) Remove() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			err = i.incomeSchedulerService.Remove(id)

			rest.PerformResponseWithCode(writer, nil, http.StatusNoContent, err)
		}
	}
}

func (i *IncomeSchedulerHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			response, err := i.incomeSchedulerService.FindById(id)

			rest.PerformResponse(writer, response, err)
		}
	}
}

func (i *IncomeSchedulerHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			response := i.incomeSchedulerService.FindByHouseId(id)

			rest.PerformResponse(writer, response, nil)
		}
	}
}
