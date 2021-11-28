// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	model "github.com/VlasovArtem/hob/src/user/model"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: user
func (_m *UserRepository) Create(user model.User) (model.User, error) {
	ret := _m.Called(user)

	var r0 model.User
	if rf, ok := ret.Get(0).(func(model.User) model.User); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExistsByEmail provides a mock function with given fields: email
func (_m *UserRepository) ExistsByEmail(email string) bool {
	ret := _m.Called(email)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ExistsById provides a mock function with given fields: id
func (_m *UserRepository) ExistsById(id uuid.UUID) bool {
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
func (_m *UserRepository) FindById(id uuid.UUID) (model.User, error) {
	ret := _m.Called(id)

	var r0 model.User
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.User); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
