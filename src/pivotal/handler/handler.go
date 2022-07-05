package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	pivotalService "github.com/VlasovArtem/hob/src/pivotal/service"
	"github.com/gorilla/mux"
	"net/http"
)

type PivotalHandlerObject struct {
	pivotalService pivotalService.PivotalService
}

func NewPivotalHandler(pivotalService pivotalService.PivotalService) PivotalHandler {
	return &PivotalHandlerObject{
		pivotalService: pivotalService,
	}
}

func (p *PivotalHandlerObject) Initialize(factory dependency.DependenciesProvider) any {
	return NewPivotalHandler(dependency.FindRequiredDependency[pivotalService.PivotalServiceStr, pivotalService.PivotalService](factory))
}

func (p *PivotalHandlerObject) Init(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1/pivotal").Subrouter()

	subrouter.Path("/house/{id}").HandlerFunc(p.FindByHouseId()).Methods("GET")
	subrouter.Path("/group/{id}").HandlerFunc(p.FindByHouseId()).Methods("GET")
}

type PivotalHandler interface {
	FindByHouseId() http.HandlerFunc
	FindByGroupId() http.HandlerFunc
}

func (p *PivotalHandlerObject) FindByHouseId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//if id, err := rest.GetIdRequestParameter(request); err != nil {
		//	rest.HandleWithError(writer, err)
		//} else {
		//	rest.NewAPIResponse(writer).
		//		Ok(p.pivotalService.FindByHouseId(id)).
		//		Perform()
		//}
	}
}

func (p *PivotalHandlerObject) FindByGroupId() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//if id, err := rest.GetIdRequestParameter(request); err != nil {
		//	rest.HandleWithError(writer, err)
		//} else {
		//	rest.NewAPIResponse(writer).
		//		Ok(p.pivotalService.FindByGroupId(id)).
		//		Perform()
		//}
	}
}
