package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

type APIResponse struct {
	writer     http.ResponseWriter
	body       any
	err        error
	statusCode int
}

func NewAPIResponse(writer http.ResponseWriter) *APIResponse {
	return &APIResponse{
		writer:     writer,
		statusCode: http.StatusOK,
	}
}

func (a *APIResponse) Body(body any) *APIResponse {
	a.body = body
	return a
}

func (a *APIResponse) Error(err error) *APIResponse {
	a.err = err
	return a
}

func (a *APIResponse) StatusCode(statusCode int) *APIResponse {
	a.statusCode = statusCode
	return a
}

func (a *APIResponse) Ok(body any, err error) *APIResponse {
	return a.Body(body).
		Error(err)
}

func (a *APIResponse) NoContent(err error) *APIResponse {
	return a.
		Error(err).
		StatusCode(http.StatusNoContent)
}

func (a *APIResponse) Created(body any, err error) *APIResponse {
	return a.StatusCode(http.StatusCreated).
		Error(err).
		Body(body)
}

func (a *APIResponse) Perform() {
	PerformResponseWithCode(a.writer, a.body, a.statusCode, a.err)
}

func ReadRequestBody[T any](r *http.Request) (requestBody T, err error) {
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		return requestBody, err
	} else if string(body) == "null" {
		return requestBody, errors.New("body not found")
	} else {
		if err = json.Unmarshal(body, &requestBody); err != nil {
			return requestBody, err
		}
	}

	return requestBody, err
}

func GetRequestParameter(request *http.Request, name string) (string, error) {
	vars := mux.Vars(request)

	if result, ok := vars[name]; ok {
		return result, nil
	}
	return "", errors.New(fmt.Sprintf("parameter '%s' not found", name))
}

func GetQueryIntParameterOrDefault(request *http.Request, name string, defaultValue int) (int, error) {
	parameter := request.URL.Query().Get(name)

	if parameter == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(parameter)
}

func GetIdRequestParameter(request *http.Request) (id uuid.UUID, err error) {
	if parameter, err := GetRequestParameter(request, "id"); err != nil {
		return id, err
	} else {

		id, err = uuid.Parse(parameter)

		if err != nil {
			return id, errors.New(fmt.Sprintf("the id is not valid %s", parameter))
		}

		return id, err
	}
}

func PerformResponse(writer http.ResponseWriter, body any, err error) {
	if err != nil {
		HandleWithError(writer, err)
	} else if body != nil {
		if err = json.NewEncoder(writer).Encode(body); err != nil {
			HandleErrorResponseWithError(writer, http.StatusInternalServerError, err)
		}
	}
}

func PerformResponseWithCode(writer http.ResponseWriter, body any, statusCode int, err error) {
	if err != nil {
		HandleWithError(writer, err)
	} else {
		writer.WriteHeader(statusCode)

		if body != nil {
			if err = json.NewEncoder(writer).Encode(body); err != nil {
				HandleErrorResponseWithError(writer, http.StatusInternalServerError, err)
			}
		}
	}
}

func HandleErrorResponseWithError(writer http.ResponseWriter, statusCode int, err error) {
	message := err.Error()

	http.Error(writer, message, statusCode)
}

func HandleWithError(writer http.ResponseWriter, err error) {
	if errors.Is(err, int_errors.ErrNotFound{}) {
		HandleErrorResponseWithError(writer, http.StatusNotFound, err)
	} else {
		HandleErrorResponseWithError(writer, http.StatusBadRequest, err)
	}
}

func HandleBadRequestWithErrorResponse(writer http.ResponseWriter, response int_errors.ErrorResponse) bool {
	if response != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(response)
		return true
	}

	return false
}
