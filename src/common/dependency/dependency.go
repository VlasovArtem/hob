package dependency

import (
	"github.com/rs/zerolog/log"
	"reflect"
)

type dependenciesProviderObject struct {
	dependenciesByName map[string]interface{}
	dependenciesByType map[reflect.Type]interface{}
}

func NewDependenciesProvider() DependenciesProvider {
	return &dependenciesProviderObject{make(map[string]interface{}), make(map[reflect.Type]interface{})}
}

type DependenciesProvider interface {
	Add(interface{}) interface{}
	AddAutoDependency(initializer ObjectDependencyInitializer) ObjectDependencyInitializer
	FindByName(dependencyName string, required bool) interface{}
	FindByType(typeOf reflect.Type, required bool) interface{}
	FindRequiredByType(typeOf reflect.Type) interface{}
	FindByObject(object interface{}) interface{}
	FindRequired(dependencyName string) interface{}
	FindRequiredByObject(object interface{}) interface{}
}

type ObjectDependencyInitializer interface {
	Initialize(factory DependenciesProvider) interface{}
}

type ObjectDatabaseMigrator interface {
	ObjectDependencyInitializer
	GetEntity() interface{}
}

func (d *dependenciesProviderObject) Add(dependency interface{}) interface{} {
	name, typeOf := findNameAndType(dependency)

	if _, exists := d.dependenciesByName[name]; !exists {
		log.Info().Msgf("dependency with name %s added", name)
		d.dependenciesByName[name] = dependency
		d.dependenciesByType[typeOf] = dependency
		return dependency
	}
	log.Info().Msgf("dependency with name %s already exists", name)
	return dependency
}

func (d *dependenciesProviderObject) FindByName(dependencyName string, required bool) interface{} {
	dependency := d.dependenciesByName[dependencyName]
	if required && dependency == nil {
		log.Fatal().Msgf("dependency with name %s not found", dependencyName)
	}
	return dependency
}

func (d *dependenciesProviderObject) FindByType(typeOf reflect.Type, required bool) interface{} {
	dependency := d.dependenciesByType[typeOf]
	if required && dependency == nil {
		log.Fatal().Msgf("dependency with type %s not found", typeOf)
	}
	return dependency
}

func (d *dependenciesProviderObject) FindRequiredByType(typeOf reflect.Type) interface{} {
	return d.FindByType(typeOf, true)
}

func (d *dependenciesProviderObject) FindByObject(object interface{}) interface{} {
	return d.FindByName(findName(object), false)
}

func (d *dependenciesProviderObject) FindRequired(dependencyName string) interface{} {
	return d.FindByName(dependencyName, true)
}

func (d *dependenciesProviderObject) FindRequiredByObject(object interface{}) interface{} {
	return d.FindRequired(findName(object))
}

func (d *dependenciesProviderObject) AddAutoDependency(initializer ObjectDependencyInitializer) ObjectDependencyInitializer {
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

func findNameAndType(dependency interface{}) (name string, dependencyType reflect.Type) {
	valueOf := reflect.Indirect(reflect.ValueOf(dependency))
	typeOf := valueOf.Type()

	switch typeOf.Kind() {
	case reflect.Struct:
		name = typeOf.Name()
	case reflect.Ptr:
		name = typeOf.Elem().Name()
	default:
		log.Fatal().Msgf("type is supported %s", typeOf.Kind())
	}

	return name, typeOf
}
