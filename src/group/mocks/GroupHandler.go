// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// GroupHandler is an autogenerated mock type for the GroupHandler type
type GroupHandler struct {
	mock.Mock
}

// Add provides a mock function with given fields:
func (_m *GroupHandler) Add() http.HandlerFunc {
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
func (_m *GroupHandler) Delete() http.HandlerFunc {
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
func (_m *GroupHandler) FindById() http.HandlerFunc {
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

// FindByUserId provides a mock function with given fields:
func (_m *GroupHandler) FindByUserId() http.HandlerFunc {
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
func (_m *GroupHandler) Update() http.HandlerFunc {
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
