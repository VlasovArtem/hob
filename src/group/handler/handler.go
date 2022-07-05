package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/group/model"
	"github.com/VlasovArtem/hob/src/group/service"
	"github.com/gorilla/mux"
	"net/http"
)

type GroupHandlerObject struct {
	groupService service.GroupService
}

func NewGroupHandler(groupService service.GroupService) GroupHandler {
	return &GroupHandlerObject{groupService}
}

func (g *GroupHandlerObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewGroupHandler(dependency.FindRequiredDependency[service.GroupServiceStr, service.GroupService](factory))
}

func (g *GroupHandlerObject) Init(router *mux.Router) {
	incomeRouter := router.PathPrefix("/api/v1/groups").Subrouter()

	incomeRouter.Path("").HandlerFunc(g.Add()).Methods("POST")
	incomeRouter.Path("/batch").HandlerFunc(g.AddBatch()).Methods("POST")
	incomeRouter.Path("/{id}").HandlerFunc(g.FindById()).Methods("GET")
	incomeRouter.Path("/user/{id}").HandlerFunc(g.FindByUserId()).Methods("GET")
	incomeRouter.Path("/{id}").HandlerFunc(g.Delete()).Methods("DELETE")
	incomeRouter.Path("/{id}").HandlerFunc(g.Update()).Methods("PUT")
}

type GroupHandler interface {
	Add() http.HandlerFunc
	AddBatch() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByUserId() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
}

func (g *GroupHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if body, err := rest.ReadRequestBody[model.CreateGroupRequest](request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Created(g.groupService.Add(body)).
				Perform()
		}
	}
}

func (g *GroupHandlerObject) AddBatch() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if body, err := rest.ReadRequestBody[model.CreateGroupBatchRequest](request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Created(g.groupService.AddBatch(body)).
				Perform()
		}
	}
}

func (g *GroupHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Ok(g.groupService.FindById(id)).
				Perform()
		}
	}
}

func (g *GroupHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Body(g.groupService.FindByUserId(id)).
				Perform()
		}
	}
}

func (g *GroupHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if body, err := rest.ReadRequestBody[model.UpdateGroupRequest](request); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				rest.NewAPIResponse(writer).
					Error(g.groupService.Update(id, body)).
					Perform()
			}
		}
	}
}

func (g *GroupHandlerObject) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				StatusCode(http.StatusNoContent).
				Error(g.groupService.DeleteById(id)).
				Perform()
		}
	}
}
