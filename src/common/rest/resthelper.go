package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/int-errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

var nilTime *time.Time

var mappers = map[reflect.Type]func(value string) (any, error){
	reflect.TypeOf(uuid.UUID{}): func(value string) (any, error) {
		if parse, err := uuid.Parse(value); err != nil {
			return uuid.UUID{}, errors.New("the id is not valid UUID")
		} else {
			return parse, nil
		}
	},
	reflect.TypeOf(int64(0)): func(value string) (any, error) {
		return strconv.ParseInt(value, 10, 64)
	},
	reflect.TypeOf(0): func(value string) (any, error) {
		return strconv.Atoi(value)
	},
	reflect.TypeOf(""): func(value string) (any, error) {
		return value, nil
	},
	reflect.TypeOf(nilTime): func(value string) (any, error) {
		if parse, err := time.Parse(time.RFC3339, value); err != nil {
			return nil, errors.New("the time is not valid RFC3339")
		} else {
			return &parse, nil
		}
	},
}

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

func GetQueryParamOrDefault[T any](request *http.Request, name string, defaultValue T) (T, error) {
	parameter := request.URL.Query().Get(name)

	if parameter == "" {
		return defaultValue, nil
	}

	var t T

	if mapper, ok := mappers[reflect.TypeOf(t)]; ok {
		if mappedValue, err := mapper(parameter); err != nil {
			return t, err
		} else {
			return mappedValue.(T), nil
		}
	} else {
		return t, errors.New(fmt.Sprintf("mapper not found %s", reflect.TypeOf(t)))
	}
}

func GetQueryParamOrDefaultReference[T any](request *http.Request, name string, defaultValue *T) (*T, error) {
	parameter := request.URL.Query().Get(name)

	if parameter == "" {
		return defaultValue, nil
	}

	var t *T

	if mapper, ok := mappers[reflect.TypeOf(t)]; ok {
		if mappedValue, err := mapper(parameter); err != nil {
			return nil, err
		} else {
			return mappedValue.(*T), nil
		}
	} else {
		return nil, errors.New(fmt.Sprintf("mapper not found %s", reflect.TypeOf(t)))
	}
}

func GetQueryParam[T any](request *http.Request, name string) (T, error) {
	parameter := request.URL.Query().Get(name)

	var t T

	if parameter == "" {
		return t, errors.New(fmt.Sprintf("parameter '%s' not found", name))
	}

	if mapper, ok := mappers[reflect.TypeOf(t)]; ok {
		if mappedValue, err := mapper(parameter); err != nil {
			return t, err
		} else {
			return mappedValue.(T), nil
		}
	} else {
		return t, errors.New(fmt.Sprintf("mapper not found %s", reflect.TypeOf(t)))
	}
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
	} else if errors.Is(err, int_errors.ErrResponse{}) {
		handleBadRequestWithErrorResponse(writer, err.(*int_errors.ErrResponse).Response)
	} else {
		HandleErrorResponseWithError(writer, http.StatusBadRequest, err)
	}
}

func handleBadRequestWithErrorResponse(writer http.ResponseWriter, response int_errors.ErrorResponse) {
	if response != nil {
		writer.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(writer).Encode(response); err != nil {
			log.Error().Err(err).Msg("ErrorResponse encoding failure")
		}
	}
}

func GetRequestPaging(request *http.Request, defaultLimit, defaultOffset int) (limit, offset int) {
	limit, err := GetQueryParamOrDefault(request, "limit", defaultLimit)

	if err != nil {
		log.Err(err).Msg("Failed to get limit query parameter")
		limit = defaultLimit
	}

	offset, err = GetQueryParamOrDefault(request, "offset", defaultOffset)

	if err != nil {
		log.Err(err).Msg("Failed to get offset query parameter")
		offset = defaultOffset
	}

	return limit, offset
}

func GetRequestFiltering(request *http.Request) (from, to *time.Time) {
	from, err := GetQueryParamOrDefaultReference[time.Time](request, "from", nil)
	if err != nil {
		log.Err(err).Msg("Failed to get from query parameter 'from'")
		return nil, nil
	}

	to, err = GetQueryParamOrDefaultReference[time.Time](request, "to", nil)
	if err != nil {
		log.Err(err).Msg("Failed to get to query parameter 'to'")
		return from, nil
	}

	return from, to
}
