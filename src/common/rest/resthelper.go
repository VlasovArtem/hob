package rest

import (
	projectErrors "common/errors"
	"common/service"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func PerformRequest(target interface{}, w http.ResponseWriter, r *http.Request) error {
	reqBody, err := ioutil.ReadAll(r.Body)

	service.LogError(err, "")

	err = json.Unmarshal(reqBody, &target)

	if target == nil {
		err = errors.New("request not parsed")
	}

	if service.LogError(err, "") {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}

	return err
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

func GetIdRequestParameter(request *http.Request) (uuid.UUID, error) {
	if parameter, err := GetRequestParameter(request, "id"); err != nil {
		return [16]byte{}, err
	} else {

		id, err := uuid.Parse(parameter)

		if err != nil {
			return [16]byte{}, errors.New(fmt.Sprintf("the id is not valid %s", parameter))
		}

		return id, nil
	}
}

func PerformResponse(writer http.ResponseWriter, response interface{}, err error) {
	if err != nil {
		HandleBadRequestWithError(writer, err)
	} else if response != nil {
		if err = json.NewEncoder(writer).Encode(response); err != nil {
			HandleErrorResponseWithError(writer, http.StatusInternalServerError, err)
		}
	}
}

func PerformResponseWithCode(writer http.ResponseWriter, response interface{}, statusCode int, err error) {
	if err != nil {
		HandleBadRequestWithError(writer, err)
	} else {
		writer.WriteHeader(statusCode)

		if response != nil {
			if err = json.NewEncoder(writer).Encode(response); err != nil {
				HandleErrorResponseWithError(writer, http.StatusInternalServerError, err)
			}
		}
	}
}

func HandleErrorResponseWithError(writer http.ResponseWriter, statusCode int, err error) {
	message := err.Error()

	http.Error(writer, message, statusCode)
}

func HandleBadRequestWithError(writer http.ResponseWriter, err error) {
	HandleErrorResponseWithError(writer, http.StatusBadRequest, err)
}

func HandleBadRequestWithErrorResponse(writer http.ResponseWriter, response projectErrors.ErrorResponse) bool {
	if response != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(response)
		return true
	}

	return false
}
