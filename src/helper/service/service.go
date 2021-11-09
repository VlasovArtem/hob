package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func PerformRequest(target interface{}, w http.ResponseWriter, r *http.Request) error {
	reqBody, err := ioutil.ReadAll(r.Body)

	HandleError(err, "")

	err = json.Unmarshal(reqBody, &target)

	if HandleError(err, "") {
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
	return "", errors.New("parameter not found")
}

func HandleErrorResponse(writer http.ResponseWriter, statusCode int, message string) {
	http.Error(writer, message, statusCode)
	writer.Write([]byte(message))
}

func HandleErrorResponseWithError(writer http.ResponseWriter, statusCode int, err error) {
	message := err.Error()

	http.Error(writer, message, statusCode)
	writer.Write([]byte(message))
}

func HandleBadRequestWithError(writer http.ResponseWriter, err error) {
	HandleErrorResponseWithError(writer, http.StatusBadRequest, err)
}

func HandleError(err error, message string) bool {
	if err != nil {
		log.Println(fmt.Sprint(err, message))
		return true
	}
	return false
}
