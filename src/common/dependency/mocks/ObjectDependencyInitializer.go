// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import (
	dependency "github.com/VlasovArtem/hob/src/common/dependency"
	mock "github.com/stretchr/testify/mock"
)

// ObjectDependencyInitializer is an autogenerated mock type for the ObjectDependencyInitializer type
type ObjectDependencyInitializer struct {
	mock.Mock
}

// Initialize provides a mock function with given fields: factory
func (_m *ObjectDependencyInitializer) Initialize(factory dependency.DependenciesProvider) interface{} {
	ret := _m.Called(factory)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(dependency.DependenciesProvider) interface{}); ok {
		r0 = rf(factory)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}
