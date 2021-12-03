package dependency

import (
	"github.com/rs/zerolog/log"
	"reflect"
)

var initializerType = reflect.TypeOf((*ObjectDependencyInitializer)(nil)).Elem()

type dependencyFactoryObject struct {
	dependencies map[string]interface{}
}

func NewDependenciesFactory() DependenciesFactory {
	return &dependencyFactoryObject{make(map[string]interface{})}
}

type DependenciesFactory interface {
	Add(interface{}) interface{}
	AddAutoDependency(initializer ObjectDependencyInitializer) ObjectDependencyInitializer
	Find(dependencyName string, required bool) interface{}
	FindByObject(object interface{}) interface{}
	FindRequired(dependencyName string) interface{}
	FindRequiredByObject(object interface{}) interface{}
}

type ObjectDependencyInitializer interface {
	Initialize(factory DependenciesFactory) interface{}
}

type ObjectDatabaseMigrator interface {
	ObjectDependencyInitializer
	GetEntity() interface{}
}

func (d *dependencyFactoryObject) Add(dependency interface{}) interface{} {
	name := findName(dependency)

	if _, exists := d.dependencies[name]; !exists {
		log.Info().Msgf("dependency with name %s added", name)
		d.dependencies[name] = dependency
		return dependency
	}
	log.Info().Msgf("dependency with name %s already exists", name)
	return dependency
}

func (d *dependencyFactoryObject) Find(dependencyName string, required bool) interface{} {
	dependency := d.dependencies[dependencyName]
	if required && dependency == nil {
		log.Fatal().Msgf("dependency with name %s not found", dependencyName)
	}
	return dependency
}

func (d *dependencyFactoryObject) FindByObject(object interface{}) interface{} {
	return d.Find(findName(object), false)
}

func (d *dependencyFactoryObject) FindRequired(dependencyName string) interface{} {
	return d.Find(dependencyName, true)
}

func (d *dependencyFactoryObject) FindRequiredByObject(object interface{}) interface{} {
	return d.FindRequired(findName(object))
}

func (d *dependencyFactoryObject) AddAutoDependency(initializer ObjectDependencyInitializer) ObjectDependencyInitializer {
	return d.Add(initializer.Initialize(d)).(ObjectDependencyInitializer)
}

func findName(dependency interface{}) (name string) {
	typeOf := reflect.TypeOf(dependency)

	switch typeOf.Kind() {
	case reflect.Struct:
		name = typeOf.Name()
	case reflect.Ptr:
		name = typeOf.Elem().Name()
	default:
		log.Fatal().Msgf("type is supported %s", typeOf.Kind())
	}

	return name
}
