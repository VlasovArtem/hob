// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	model "user/model"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// Add provides a mock function with given fields: request
func (_m *UserService) Add(request model.CreateUserRequest) (model.UserResponse, error) {
	ret := _m.Called(request)

	var r0 model.UserResponse
	if rf, ok := ret.Get(0).(func(model.CreateUserRequest) model.UserResponse); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Get(0).(model.UserResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.CreateUserRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExistsById provides a mock function with given fields: id
func (_m *UserService) ExistsById(id uuid.UUID) bool {
	ret := _m.Called(id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uuid.UUID) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FindById provides a mock function with given fields: id
func (_m *UserService) FindById(id uuid.UUID) (model.UserResponse, error) {
	ret := _m.Called(id)

	var r0 model.UserResponse
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.UserResponse); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.UserResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
