package handler

import (
	"encoding/json"
	"github.com/VlasovArtem/hob/src/common/dependency"
	projectErrors "github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/VlasovArtem/hob/src/common/rest"
	"github.com/VlasovArtem/hob/src/user/model"
	"github.com/VlasovArtem/hob/src/user/service"
	"github.com/VlasovArtem/hob/src/user/validator"
	"github.com/google/uuid"
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

func (u *UserHandlerObject) Initialize(factory dependency.DependenciesProvider) interface{} {
	return NewUserHandler(
		factory.FindRequiredByObject(service.UserServiceObject{}).(service.UserService),
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
		requestEntity := model.CreateUserRequest{}

		if err := rest.PerformRequest(&requestEntity, writer, request); err != nil {
			rest.HandleWithError(writer, err)
			return
		}
		if rest.HandleBadRequestWithErrorResponse(writer, u.userValidator.ValidateCreateRequest(requestEntity)) {
			return
		}

		if userResponse, err := u.userService.Add(requestEntity); err != nil {
			rest.HandleBadRequestWithErrorResponse(writer, projectErrors.NewWithDetails(err.Error()))
		} else {
			writer.WriteHeader(http.StatusCreated)
			json.NewEncoder(writer).Encode(userResponse)
		}
	}
}

func (u *UserHandlerObject) FindById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if parameter, err := rest.GetRequestParameter(request, "id"); err != nil {
			rest.HandleWithError(writer, err)

			return
		} else {
			if id, err := uuid.Parse(parameter); err != nil {
				rest.HandleWithError(writer, err)
			} else {
				userResponse, err := u.userService.FindById(id)
				rest.PerformResponse(writer, userResponse, err)
			}
		}
	}
}

func (u *UserHandlerObject) Delete() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			rest.PerformResponseWithCode(writer, nil, http.StatusNoContent, u.userService.Delete(id))
		}
	}
}

func (u *UserHandlerObject) Update() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if id, err := rest.GetIdRequestParameter(request); err != nil {
			rest.HandleWithError(writer, err)
		} else {
			requestEntity := model.UpdateUserRequest{}

			if err := rest.PerformRequest(&requestEntity, writer, request); err != nil {
				rest.HandleWithError(writer, err)
				return
			}

			if rest.HandleBadRequestWithErrorResponse(writer, u.userValidator.ValidateUpdateRequest(requestEntity)) {
				return
			}

			if err = u.userService.Update(id, requestEntity); err != nil {
				rest.HandleBadRequestWithErrorResponse(writer, projectErrors.NewWithDetails(err.Error()))
			} else {
				rest.PerformResponse(writer, nil, nil)
			}
		}
	}
}
