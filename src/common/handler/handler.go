package handler

import (
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/gorilla/mux"
)

type ApplicationHandler interface {
	dependency.ObjectDependencyInitializer
	Init(*mux.Router)
}
