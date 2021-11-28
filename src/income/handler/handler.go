package handler

import (
	"common/dependency"
	"common/rest"
	"github.com/gorilla/mux"
	"income/model"
	"income/service"
	"net/http"
)

type IncomeHandlerObject struct {
	incomeService service.IncomeService
}

func NewIncomeHandler(incomeService service.IncomeService) IncomeHandler {
	return &IncomeHandlerObject{incomeService}
}

func (i *IncomeHandlerObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(NewIncomeHandler(factory.FindRequiredByObject(service.IncomeServiceObject{}).(service.IncomeService)))
}

func (i *IncomeHandlerObject) Init(router *mux.Router) {
	incomeRouter := router.PathPrefix("/api/v1/income").Subrouter()

	incomeRouter.Path("").HandlerFunc(i.Add()).Methods("POST")
	incomeRouter.Path("/{id}").HandlerFunc(i.FindById()).Methods("GET")
	incomeRouter.Path("/house/{id}").HandlerFunc(i.FindByHouseId()).Methods("GET")
}

type IncomeHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
}

func (i *IncomeHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var requestBody model.CreateIncomeRequest

		if err := rest.PerformRequest(&requestBody, writer, request); err != nil {
			return
		}

		meter, err := i.incomeService.Add(requestBody)

		rest.PerformResponseWithCode(writer, meter, http.StatusCreated, err)
	}
}

func (i *IncomeHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			meterResponse, err := i.incomeService.FindById(id)

			rest.PerformResponse(writer, meterResponse, err)
		}
	}
}

func (i *IncomeHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
		} else {
			rest.PerformResponse(writer, i.incomeService.FindByHouseId(id), err)
		}
	}
}
