package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/income/model"
	"github.com/VlasovArtem/hob/src/income/service"
	"github.com/gorilla/mux"
	"net/http"
)

type IncomeHandlerStr struct {
	incomeService service.IncomeService
}

func NewIncomeHandler(incomeService service.IncomeService) IncomeHandler {
	return &IncomeHandlerStr{incomeService}
}

func (i *IncomeHandlerStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(service.IncomeServiceStr{}),
	}
}

func (i *IncomeHandlerStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewIncomeHandler(dependency.FindRequiredDependency[service.IncomeServiceStr, service.IncomeService](factory))
}

func (i *IncomeHandlerStr) Init(router *mux.Router) {
	incomeRouter := router.PathPrefix("/api/v1/incomes").Subrouter()

	incomeRouter.Path("").HandlerFunc(i.Add()).Methods("POST")
	incomeRouter.Path("/batch").HandlerFunc(i.AddBatch()).Methods("POST")
	incomeRouter.Path("/{id}").HandlerFunc(i.FindById()).Methods("GET")
	incomeRouter.Path("/{id}").HandlerFunc(i.Delete()).Methods("DELETE")
	incomeRouter.Path("/{id}").HandlerFunc(i.Update()).Methods("PUT")
	incomeRouter.Path("/house/{id}").HandlerFunc(i.FindByHouseId()).Methods("GET")
}

type IncomeHandler interface {
	Add() http.HandlerFunc
	Delete() http.HandlerFunc
	Update() http.HandlerFunc
	AddBatch() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
}

func (i *IncomeHandlerStr) Add() http.HandlerFunc {
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

func (i *IncomeHandlerStr) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				StatusCode(http.StatusNoContent).
				Error(i.incomeService.DeleteById(id)).
				Perform()
		}
	}
}

func (i *IncomeHandlerStr) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if body, err := rest.ReadRequestBody[model.UpdateIncomeRequest](request); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				rest.NewAPIResponse(writer).
					Error(i.incomeService.Update(id, body)).
					Perform()
			}
		}
	}
}

func (i *IncomeHandlerStr) AddBatch() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if body, err := rest.ReadRequestBody[model.CreateIncomeBatchRequest](request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Created(i.incomeService.AddBatch(body)).
				Perform()
		}
	}
}

func (i *IncomeHandlerStr) FindById() http.HandlerFunc {
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

func (i *IncomeHandlerStr) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			limit, offset := rest.GetRequestPaging(request, 25, 0)
			from, to := rest.GetRequestFiltering(request)

			rest.NewAPIResponse(writer).
				Body(i.incomeService.FindByHouseId(id, limit, offset, from, to)).
				Perform()
		}
	}
}
