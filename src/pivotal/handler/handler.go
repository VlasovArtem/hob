package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/VlasovArtem/hob/src/pivotal/service"
	"github.com/gorilla/mux"
	"net/http"
)

type PivotalHandlerStr struct {
	pivotalService service.PivotalService
}

func NewPivotalHandler(pivotalService service.PivotalService) PivotalHandler {
	return &PivotalHandlerStr{
		pivotalService: pivotalService,
	}
}

func (p *PivotalHandlerStr) GetRequiredDependencies() []dependency.Requirements {
	return []dependency.Requirements{
		dependency.FindNameAndType(service.PivotalServiceStr{}),
	}
}

func (p *PivotalHandlerStr) Initialize(factory dependency.DependenciesProvider) any {
	return NewPivotalHandler(dependency.FindRequiredDependency[service.PivotalServiceStr, service.PivotalService](factory))
}

func (p *PivotalHandlerStr) Init(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1/pivotal").Subrouter()

	subrouter.Path("/house/{id}").HandlerFunc(p.FindByHouseId()).Methods("GET")
	subrouter.Path("/group/{id}").HandlerFunc(p.FindByHouseId()).Methods("GET")
}

type PivotalHandler interface {
	FindByHouseId() http.HandlerFunc
	FindByGroupId() http.HandlerFunc
}

func (p *PivotalHandlerStr) FindByHouseId() http.HandlerFunc {
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

func (p *PivotalHandlerStr) FindByGroupId() http.HandlerFunc {
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
