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

func (i *IncomeSchedulerHandlerObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewIncomeSchedulerHandler(dependency.FindRequiredDependency[service.IncomeSchedulerServiceStr, service.IncomeSchedulerService](factory))
}

func (i *IncomeSchedulerHandlerObject) Init(router *mux.Router) {
	incomeSchedulerRouter := router.PathPrefix("/api/v1/incomes/schedulers").Subrouter()

	incomeSchedulerRouter.Path("").HandlerFunc(i.Add()).Methods("POST")
	incomeSchedulerRouter.Path("/{id}").HandlerFunc(i.FindById()).Methods("GET")
	incomeSchedulerRouter.Path("/{id}").HandlerFunc(i.Update()).Methods("PUT")
	incomeSchedulerRouter.Path("/{id}").HandlerFunc(i.Remove()).Methods("DELETE")
	incomeSchedulerRouter.Path("/house/{id}").HandlerFunc(i.FindByHouseId()).Methods("GET")
}

type IncomeSchedulerHandler interface {
	Add() http.HandlerFunc
	Update() http.HandlerFunc
	Remove() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByHouseId() http.HandlerFunc
}

func (i *IncomeSchedulerHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if body, err := rest.ReadRequestBody[model.CreateIncomeSchedulerRequest](request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Created(i.incomeSchedulerService.Add(body)).
				Perform()
		}
	}
}

func (i *IncomeSchedulerHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if body, err := rest.ReadRequestBody[model.UpdateIncomeSchedulerRequest](request); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				rest.NewAPIResponse(writer).
					Error(i.incomeSchedulerService.Update(id, body)).
					Perform()
			}
		}
	}
}

func (i *IncomeSchedulerHandlerObject) Remove() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				StatusCode(http.StatusNoContent).
				Error(i.incomeSchedulerService.DeleteById(id)).
				Perform()
		}
	}
}

func (i *IncomeSchedulerHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Ok(i.incomeSchedulerService.FindById(id)).
				Perform()
		}
	}
}

func (i *IncomeSchedulerHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Body(i.incomeSchedulerService.FindByHouseId(id)).
				Perform()
		}
	}
}
