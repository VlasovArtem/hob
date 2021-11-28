package handler

import (
	"common/dependency"
	projectErrors "common/errors"
	"common/rest"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"user/model"
	"user/service"
	"user/validator"
)

type UserHandlerObject struct {
	userService   service.UserService
	userValidator validator.UserRequestValidator
}

func NewUserHandler(userService service.UserService, userValidator validator.UserRequestValidator) UserHandler {
	return &UserHandlerObject{userService, userValidator}
}

func (u *UserHandlerObject) Initialize(factory dependency.DependenciesFactory) {
	factory.Add(
		NewUserHandler(
			factory.FindRequiredByObject(service.UserServiceObject{}).(service.UserService),
			factory.FindRequiredByObject(validator.UserRequestValidatorObject{}).(validator.UserRequestValidator),
		),
	)
}

type UserHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
}

func (u *UserHandlerObject) Init(router *mux.Router) {
	userRouter := router.PathPrefix("/api/v1/user").Subrouter()

	userRouter.Path("").HandlerFunc(u.Add()).Methods("POST")
	userRouter.Path("/{id}").HandlerFunc(u.FindById()).Methods("GET")
}

func (u *UserHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestEntity := model.CreateUserRequest{}

		if err := rest.PerformRequest(&requestEntity, writer, request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
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
			rest.HandleBadRequestWithError(writer, err)

			return
		} else {
			if id, err := uuid.Parse(parameter); err != nil {
				rest.HandleBadRequestWithError(writer, err)
			} else {
				if userResponse, err := u.userService.FindById(id); err != nil {
					rest.HandleErrorResponseWithError(writer, http.StatusNotFound, err)
				} else {
					json.NewEncoder(writer).Encode(userResponse)
				}
			}
		}
	}
}
