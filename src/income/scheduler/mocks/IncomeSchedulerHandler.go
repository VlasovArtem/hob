// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// IncomeSchedulerHandler is an autogenerated mock type for the IncomeSchedulerHandler type
type IncomeSchedulerHandler struct {
	mock.Mock
}

// Add provides a mock function with given fields:
func (_m *IncomeSchedulerHandler) Add() http.HandlerFunc {
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

// FindByHouseId provides a mock function with given fields:
func (_m *IncomeSchedulerHandler) FindByHouseId() http.HandlerFunc {
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
func (_m *IncomeSchedulerHandler) FindById() http.HandlerFunc {
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

// Remove provides a mock function with given fields:
func (_m *IncomeSchedulerHandler) Remove() http.HandlerFunc {
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
