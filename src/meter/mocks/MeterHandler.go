// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MeterHandler is an autogenerated mock type for the MeterHandler type
type MeterHandler struct {
	mock.Mock
}

// Add provides a mock function with given fields:
func (_m *MeterHandler) Add() http.HandlerFunc {
	ret := _m.Called()

	var r0 http.HandlerFunc
	if rf, ok := ret.Get(0).(func() http.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.HandlerFunc)
		}
	}

	return r0
}

// Delete provides a mock function with given fields:
func (_m *MeterHandler) Delete() http.HandlerFunc {
	ret := _m.Called()

	var r0 http.HandlerFunc
	if rf, ok := ret.Get(0).(func() http.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.HandlerFunc)
		}
	}

	return r0
}

// FindById provides a mock function with given fields:
func (_m *MeterHandler) FindById() http.HandlerFunc {
	ret := _m.Called()

	var r0 http.HandlerFunc
	if rf, ok := ret.Get(0).(func() http.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.HandlerFunc)
		}
	}

	return r0
}

// FindByPaymentId provides a mock function with given fields:
func (_m *MeterHandler) FindByPaymentId() http.HandlerFunc {
	ret := _m.Called()

	var r0 http.HandlerFunc
	if rf, ok := ret.Get(0).(func() http.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.HandlerFunc)
		}
	}

	return r0
}

// Update provides a mock function with given fields:
func (_m *MeterHandler) Update() http.HandlerFunc {
	ret := _m.Called()

	var r0 http.HandlerFunc
	if rf, ok := ret.Get(0).(func() http.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.HandlerFunc)
		}
	}

	return r0
}

type mockConstructorTestingTNewMeterHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewMeterHandler creates a new instance of MeterHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMeterHandler(t mockConstructorTestingTNewMeterHandler) *MeterHandler {
	mock := &MeterHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
