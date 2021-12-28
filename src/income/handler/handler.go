package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/income/service"
	"github.com/gorilla/mux"
	"net/http"
)

type IncomeHandlerObject struct {
	incomeService service.IncomeService
}

func NewIncomeHandler(incomeService service.IncomeService) IncomeHandler {
	return &IncomeHandlerObject{incomeService}
}

func (i *IncomeHandlerObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewIncomeHandler(factory.FindRequiredByObject(service.IncomeServiceObject{}).(service.IncomeService))
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
		if body, err := rest.ReadRequestBody[model.CreateIncomeRequest](request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Created(i.incomeService.Add(body)).
				Perform()
		}
	}
}

func (i *IncomeHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Ok(i.incomeService.FindById(id)).
				Perform()
		}
	}
}

func (i *IncomeHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Body(i.incomeService.FindByHouseId(id)).
				Perform()
		}
	}
}
