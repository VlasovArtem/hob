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

func (i *GroupHandlerObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewGroupHandler(factory.FindRequiredByObject(service.GroupServiceType).(service.GroupService))
}

func (i *GroupHandlerObject) Init(router *mux.Router) {
	incomeRouter := router.PathPrefix("/api/v1/group").Subrouter()

	incomeRouter.Path("").HandlerFunc(i.Add()).Methods("POST")
	incomeRouter.Path("/{id}").HandlerFunc(i.FindById()).Methods("GET")
	incomeRouter.Path("/user/{id}").HandlerFunc(i.FindByUserId()).Methods("GET")
}

type GroupHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
	FindByUserId() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
}

func (i *GroupHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if body, err := rest.ReadRequestBody[model.CreateGroupRequest](request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Created(i.groupService.Add(body)).
				Perform()
		}
	}
}

func (i *GroupHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Ok(i.groupService.FindById(id)).
				Perform()
		}
	}
}

func (i *GroupHandlerObject) FindByUserId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Body(i.groupService.FindByUserId(id)).
				Perform()
		}
	}
}

func (h *GroupHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if body, err := rest.ReadRequestBody[model.UpdateGroupRequest](request); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				rest.NewAPIResponse(writer).
					Error(h.groupService.Update(id, body)).
					Perform()
			}
		}
	}
}

func (h *GroupHandlerObject) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				StatusCode(http.StatusNoContent).
				Error(h.groupService.DeleteById(id)).
				Perform()
		}
	}
}
