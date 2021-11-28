package handler

import (
	"common/dependency"
	"github.com/gorilla/mux"
)

type ApplicationHandler interface {
	dependency.ObjectDependencyInitializer
	Init(*mux.Router)
}
