package testhelper

import (
	"bytes"
	"country/model"
	countries "country/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func InitCountryService() countries.CountryService {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/content/countries.json", os.Getenv("GOPATH")))

	if err != nil {
		log.Fatal(err, "countries file not fount")
	}

	var countriesContent []model.Country

	json.Unmarshal(file, &countriesContent)

	return countries.NewCountryService(countriesContent)
}

type TestRequest struct {
	Method   string
	URL      string
	Body     interface{}
	Vars     map[string]string
	Handler  http.HandlerFunc
	Request  *http.Request
	Recorder *httptest.ResponseRecorder
	build    bool
}

func NewTestRequest() *TestRequest {
	return &TestRequest{
		Recorder: httptest.NewRecorder(),
		Vars:     make(map[string]string),
	}
}

type TestRequestBuilder interface {
	WithMethod(method string) *TestRequest
	WithURL(target string) *TestRequest
	WithBody(body interface{}) *TestRequest
	WithHandler(handler http.HandlerFunc) *TestRequest
	WithVar(key string, value string) *TestRequest
	Build() *TestRequest
}

func (t *TestRequest) WithVar(key string, value string) *TestRequest {
	t.Vars[key] = value

	return t
}

func (t *TestRequest) WithMethod(method string) *TestRequest {
	t.Method = method
	return t
}
func (t *TestRequest) WithURL(URL string) *TestRequest {
	t.URL = URL
	return t
}

func (t *TestRequest) WithBody(body interface{}) *TestRequest {
	t.Body = body
	return t
}

func (t *TestRequest) WithHandler(handler http.HandlerFunc) *TestRequest {
	t.Handler = handler
	return t
}

func (t *TestRequest) Build() *TestRequest {
	body, _ := json.Marshal(t.Body)

	buffer := bytes.Buffer{}

	buffer.Write(body)

	t.Request = httptest.NewRequest(t.Method, t.URL, &buffer)

	if len(t.Vars) != 0 {
		t.Request = mux.SetURLVars(t.Request, t.Vars)
	}

	t.build = true

	return t
}

func (t *TestRequest) execute() {
	t.Handler(t.Recorder, t.Request)
}

func (t *TestRequest) Verify(test *testing.T, expectedStatusCode int, expectedBody []byte) {
	if !t.build {
		t.Build()
	}

	t.execute()

	response := t.Recorder.Result()

	assert.Equal(test, expectedStatusCode, response.StatusCode, fmt.Sprintf("Response status code should be %d but was %d", expectedStatusCode, response.StatusCode))
	assert.Equal(test, expectedBody, ReadBytes(response))
}

func (t *TestRequest) VerifyWithStatus(test *testing.T, expectedStatusCode int) []byte {
	if !t.build {
		t.Build()
	}

	t.execute()

	response := t.Recorder.Result()

	assert.Equal(test, expectedStatusCode, response.StatusCode, fmt.Sprintf("Response status code should be %d but was %d", expectedStatusCode, response.StatusCode))

	return ReadBytes(response)
}

func ReadBytes(response *http.Response) []byte {
	buffer := bytes.Buffer{}

	buffer.ReadFrom(response.Body)

	return buffer.Bytes()
}
