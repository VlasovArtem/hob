package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"helper/model"
	"helper/service"
	"io/ioutil"
	"net/http"
)

func PerformRequest(target interface{}, w http.ResponseWriter, r *http.Request) error {
	reqBody, err := ioutil.ReadAll(r.Body)

	service.HandleError(err, "")

	err = json.Unmarshal(reqBody, &target)

	if service.HandleError(err, "") {
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

func HandleErrorResponseWithError(writer http.ResponseWriter, statusCode int, err error) {
	message := err.Error()

	http.Error(writer, message, statusCode)
}

func HandleBadRequestWithError(writer http.ResponseWriter, err error) {
	HandleErrorResponseWithError(writer, http.StatusBadRequest, err)
}

func HandleIfRequiredBadRequestWithErrorResponse(writer http.ResponseWriter, response model.ErrorResponse) bool {
	if len(response.Messages) == 0 {
		return false
	}
	writer.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(writer).Encode(response)

	return true
}

func HandleBadRequestWithErrorResponse(writer http.ResponseWriter, response model.ErrorResponse) {
	writer.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(writer).Encode(response)
}
