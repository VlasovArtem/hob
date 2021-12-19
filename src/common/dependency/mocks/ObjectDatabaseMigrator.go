// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	dependency "github.com/VlasovArtem/hob/src/common/dependency"
	mock "github.com/stretchr/testify/mock"
)

// ObjectDatabaseMigrator is an autogenerated mock type for the ObjectDatabaseMigrator type
type ObjectDatabaseMigrator struct {
	mock.Mock
}

// GetEntity provides a mock function with given fields:
func (_m *ObjectDatabaseMigrator) GetEntity() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Initialize provides a mock function with given fields: factory
func (_m *ObjectDatabaseMigrator) Initialize(factory dependency.DependenciesProvider) interface{} {
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
