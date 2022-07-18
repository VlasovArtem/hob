package dependency

import (
	"github.com/rs/zerolog/log"
	"reflect"
)

type Requirements struct {
	Name string
	Type reflect.Type
}

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
	GetRequiredDependencies() []Requirements
	Initialize(factory DependenciesProvider) any
}

type ObjectDatabaseMigrator[T any] interface {
	ObjectDependencyInitializer
	EntityProvider[T]
}

func NewEntity[T any](model T) EntityProvider[T] {
	return &entity[T]{model}
}

type entity[T any] struct {
	entity T
}

type EntityProvider[T any] interface {
	GetEntity() T
}

func (e *entity[T]) GetEntity() T {
	return e.entity
}

func (d *dependenciesProviderObject) Add(dependency any) any {
	req := FindNameAndType(dependency)

	if _, exists := d.dependenciesByName[req.Name]; !exists {
		log.Info().Msgf("dependency with name %s added", req.Name)
		d.dependenciesByName[req.Name] = dependency
		d.dependenciesByType[req.Type] = dependency
		return dependency
	}
	log.Info().Msgf("dependency with name %s already exists", req.Name)
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

func FindNameAndType(dependency any) (req Requirements) {
	valueOf := reflect.Indirect(reflect.ValueOf(dependency))
	typeOf := valueOf.Type()
	req.Type = typeOf

	switch typeOf.Kind() {
	case reflect.Struct:
		req.Name = typeOf.Name()
	case reflect.Ptr:
		req.Name = typeOf.Elem().Name()
	default:
		log.Fatal().Msgf("type is supported %s", typeOf.Kind())
	}

	return
}

func FindRequiredDependency[T any, V any](d DependenciesProvider) V {
	var t T
	typeOf := reflect.TypeOf(t)

	return d.FindByType(typeOf, true).(V)
}
