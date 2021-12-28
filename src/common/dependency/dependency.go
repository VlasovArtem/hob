package dependency

import (
	"github.com/rs/zerolog/log"
	"reflect"
)

type dependenciesProviderObject struct {
	dependenciesByName map[string]any
	dependenciesByType map[reflect.Type]any
}

func NewDependenciesProvider() DependenciesProvider {
	return &dependenciesProviderObject{make(map[string]any), make(map[reflect.Type]any)}
}

type DependenciesProvider interface {
	Add(any) any
	AddAutoDependency(initializer ObjectDependencyInitializer) ObjectDependencyInitializer
	FindByName(dependencyName string, required bool) any
	FindByType(typeOf reflect.Type, required bool) any
	FindRequiredByType(typeOf reflect.Type) any
	FindByObject(object any) any
	FindRequired(dependencyName string) any
	FindRequiredByObject(object any) any
}

type ObjectDependencyInitializer interface {
	Initialize(factory DependenciesProvider) any
}

type ObjectDatabaseMigrator interface {
	ObjectDependencyInitializer
	GetEntity() any
}

func (d *dependenciesProviderObject) Add(dependency any) any {
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

func (d *dependenciesProviderObject) FindByName(dependencyName string, required bool) any {
	dependency := d.dependenciesByName[dependencyName]
	if required && dependency == nil {
		log.Fatal().Msgf("dependency with name %s not found", dependencyName)
	}
	return dependency
}

func (d *dependenciesProviderObject) FindByType(typeOf reflect.Type, required bool) any {
	dependency := d.dependenciesByType[typeOf]
	if required && dependency == nil {
		log.Fatal().Msgf("dependency with type %s not found", typeOf)
	}
	return dependency
}

func (d *dependenciesProviderObject) FindRequiredByType(typeOf reflect.Type) any {
	return d.FindByType(typeOf, true)
}

func (d *dependenciesProviderObject) FindByObject(object any) any {
	return d.FindByName(findName(object), false)
}

func (d *dependenciesProviderObject) FindRequired(dependencyName string) any {
	return d.FindByName(dependencyName, true)
}

func (d *dependenciesProviderObject) FindRequiredByObject(object any) any {
	return d.FindRequired(findName(object))
}

func (d *dependenciesProviderObject) AddAutoDependency(initializer ObjectDependencyInitializer) ObjectDependencyInitializer {
	return d.Add(initializer.Initialize(d)).(ObjectDependencyInitializer)
}

func findName(dependency any) (name string) {
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

func findNameAndType(dependency any) (name string, dependencyType reflect.Type) {
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
