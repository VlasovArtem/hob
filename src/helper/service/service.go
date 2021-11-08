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

func HandleError(err error, message string) bool {
	if err != nil {
		log.Println(fmt.Sprint(err, message))
		return true
	}
	return false
}
