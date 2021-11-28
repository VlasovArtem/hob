// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// DatabaseService is an autogenerated mock type for the DatabaseService type
type DatabaseService struct {
	mock.Mock
}

// Create provides a mock function with given fields: value
func (_m *DatabaseService) Create(value interface{}) error {
	ret := _m.Called(value)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// D provides a mock function with given fields:
func (_m *DatabaseService) D() *gorm.DB {
	ret := _m.Called()

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func() *gorm.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// DM provides a mock function with given fields: model
func (_m *DatabaseService) DM(model interface{}) *gorm.DB {
	ret := _m.Called(model)

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(interface{}) *gorm.DB); ok {
		r0 = rf(model)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// ExistsById provides a mock function with given fields: model, id
func (_m *DatabaseService) ExistsById(model interface{}, id uuid.UUID) bool {
	ret := _m.Called(model, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}, uuid.UUID) bool); ok {
		r0 = rf(model, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ExistsByQuery provides a mock function with given fields: model, query, args
func (_m *DatabaseService) ExistsByQuery(model interface{}, query interface{}, args ...interface{}) bool {
	var _ca []interface{}
	_ca = append(_ca, model, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}, interface{}, ...interface{}) bool); ok {
		r0 = rf(model, query, args...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FindById provides a mock function with given fields: receiver, id
func (_m *DatabaseService) FindById(receiver interface{}, id uuid.UUID) error {
	ret := _m.Called(receiver, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, uuid.UUID) error); ok {
		r0 = rf(receiver, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByIdModeled provides a mock function with given fields: model, receiver, id
func (_m *DatabaseService) FindByIdModeled(model interface{}, receiver interface{}, id uuid.UUID) error {
	ret := _m.Called(model, receiver, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, interface{}, uuid.UUID) error); ok {
		r0 = rf(model, receiver, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
