package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/VlasovArtem/hob/src/user/service"
	"github.com/VlasovArtem/hob/src/user/validator"
	"github.com/gorilla/mux"
	"net/http"
)

type UserHandlerObject struct {
	userService   service.UserService
	userValidator validator.UserRequestValidator
}

func NewUserHandler(userService service.UserService, userValidator validator.UserRequestValidator) UserHandler {
	return &UserHandlerObject{userService, userValidator}
}

func (u *UserHandlerObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewUserHandler(
		factory.FindRequiredByObject(service.UserServiceType).(service.UserService),
		factory.FindRequiredByObject(validator.UserRequestValidatorObject{}).(validator.UserRequestValidator),
	)
}

type UserHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
	Delete() http.HandlerFunc
	Update() http.HandlerFunc
}

func (u *UserHandlerObject) Init(router *mux.Router) {
	userRouter := router.PathPrefix("/api/v1/user").Subrouter()

	userRouter.Path("").HandlerFunc(u.Add()).Methods("POST")
	userRouter.Path("/{id}").HandlerFunc(u.FindById()).Methods("GET")
	userRouter.Path("/{id}").HandlerFunc(u.Delete()).Methods("DELETE")
	userRouter.Path("/{id}").HandlerFunc(u.Update()).Methods("PUT")
}

func (u *UserHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if body, err := rest.ReadRequestBody[model.CreateUserRequest](request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if rest.HandleBadRequestWithErrorResponse(writer, u.userValidator.ValidateCreateRequest(body)) {
				return
			}
			rest.NewAPIResponse(writer).
				Created(u.userService.Add(body)).
				Perform()
		}
	}
}

func (u *UserHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				Ok(u.userService.FindById(id)).
				Perform()
		}
	}
}

func (u *UserHandlerObject) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.NewAPIResponse(writer).
				NoContent(u.userService.Delete(id)).
				Perform()
		}
	}
}

func (u *UserHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if body, err := rest.ReadRequestBody[model.UpdateUserRequest](request); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				if rest.HandleBadRequestWithErrorResponse(writer, u.userValidator.ValidateUpdateRequest(body)) {
					return
				}

				rest.NewAPIResponse(writer).
					Error(u.userService.Update(id, body)).
					Perform()
			}
		}
	}
}
