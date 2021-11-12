package handler

import (
	helperModel "common/model"
	"common/rest"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"user/model"
	"user/service"
)

type userHandlerObject struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandlerObject{userService}
}

type UserHandler interface {
	Add() http.HandlerFunc
	FindById() http.HandlerFunc
}

func (u *userHandlerObject) Add() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestEntity := model.CreateUserRequest{}

		if err := rest.PerformRequest(&requestEntity, writer, request); err != nil {
			rest.HandleBadRequestWithError(writer, err)
			return
		}
		if rest.HandleIfRequiredBadRequestWithErrorResponse(writer, validateCreateUserRequest(requestEntity)) {
			return
		}

		if userResponse, err := u.userService.Add(requestEntity); err != nil {
			rest.HandleBadRequestWithErrorResponse(writer, helperModel.ErrorResponse{
				Error: err.Error(),
			})
		} else {
			writer.WriteHeader(http.StatusCreated)
			json.NewEncoder(writer).Encode(userResponse)
		}
	}
}

func (u *userHandlerObject) FindById() http.HandlerFunc {
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

func validateCreateUserRequest(request model.CreateUserRequest) helperModel.ErrorResponse {
	response := helperModel.ErrorResponse{
		Error: "Create User Request Validation Error",
	}

	validateStringFieldNotEmpty(&response, request.Email, "email should not be empty")
	validateByteFieldNotEmpty(&response, request.Password, "password should not be empty")

	return response
}

func validateStringFieldNotEmpty(errorsAccumulator *helperModel.ErrorResponse, value string, message string) {
	if value == "" {
		errorsAccumulator.Messages = append(errorsAccumulator.Messages, message)
	}
}

func validateByteFieldNotEmpty(errorsAccumulator *helperModel.ErrorResponse, value []byte, message string) {
	if len(value) == 0 {
		errorsAccumulator.Messages = append(errorsAccumulator.Messages, message)
	}
}
