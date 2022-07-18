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

type UserHandlerStr struct {
	userService   service.UserService
	userValidator validator.UserRequestValidator
}

func NewUserHandler(userService service.UserService, userValidator validator.UserRequestValidator) UserHandler {
	return &UserHandlerStr{userService, userValidator}
}

func (u *UserHandlerStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(service.UserServiceStr{}),
		dependency.FindNameAndType(validator.UserRequestValidatorStr{}),
	}
}

func (u *UserHandlerStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewUserHandler(
		dependency.FindRequiredDependency[service.UserServiceStr, service.UserService](factory),
		dependency.FindRequiredDependency[validator.UserRequestValidatorStr, validator.UserRequestValidator](factory),
	)
}

type UserHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
	Delete() http.HandlerFunc
	Update() http.HandlerFunc
}

func (u *UserHandlerStr) Init(router *mux.Router) {
	userRouter := router.PathPrefix("/api/v1/users").Subrouter()

	userRouter.Path("").HandlerFunc(u.Add()).Methods("POST")
	userRouter.Path("/{id}").HandlerFunc(u.FindById()).Methods("GET")
	userRouter.Path("/{id}").HandlerFunc(u.Delete()).Methods("DELETE")
	userRouter.Path("/{id}").HandlerFunc(u.Update()).Methods("PUT")
}

func (u *UserHandlerStr) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if body, err := rest.ReadRequestBody[model.CreateUserRequest](request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if err := u.userValidator.ValidateCreateRequest(body); err != nil {
				rest.HandleWithError(writer, err)
				return
			} else {
				rest.NewAPIResponse(writer).
					Created(u.userService.Add(body)).
					Perform()
			}
		}
	}
}

func (u *UserHandlerStr) FindById() http.HandlerFunc {
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

func (u *UserHandlerStr) Delete() http.HandlerFunc {
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

func (u *UserHandlerStr) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			if body, err := rest.ReadRequestBody[model.UpdateUserRequest](request); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				if err := u.userValidator.ValidateUpdateRequest(body); err != nil {
					rest.HandleWithError(writer, err)
					return
				} else {
					rest.NewAPIResponse(writer).
						Error(u.userService.Update(id, body)).
						Perform()
				}
			}
		}
	}
}
