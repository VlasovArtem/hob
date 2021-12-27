// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	model "github.com/VlasovArtem/hob/src/income/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// IncomeService is an autogenerated mock type for the IncomeService type
type IncomeService struct {
	mock.Mock
}

// Add provides a mock function with given fields: request
func (_m *IncomeService) Add(request model.CreateIncomeRequest) (model.IncomeDto, error) {
	ret := _m.Called(request)

	var r0 model.IncomeDto
	if rf, ok := ret.Get(0).(func(model.CreateIncomeRequest) model.IncomeDto); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Get(0).(model.IncomeDto)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.CreateIncomeRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteById provides a mock function with given fields: id
func (_m *IncomeService) DeleteById(id uuid.UUID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExistsById provides a mock function with given fields: id
func (_m *IncomeService) ExistsById(id uuid.UUID) bool {
	ret := _m.Called(id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uuid.UUID) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FindByHouseId provides a mock function with given fields: id
func (_m *IncomeService) FindByHouseId(id uuid.UUID) []model.IncomeDto {
	ret := _m.Called(id)

	var r0 []model.IncomeDto
	if rf, ok := ret.Get(0).(func(uuid.UUID) []model.IncomeDto); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.IncomeDto)
		}
	}

	return r0
}

// FindById provides a mock function with given fields: id
func (_m *IncomeService) FindById(id uuid.UUID) (model.IncomeDto, error) {
	ret := _m.Called(id)

	var r0 model.IncomeDto
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.IncomeDto); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.IncomeDto)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: id, request
func (_m *IncomeService) Update(id uuid.UUID, request model.UpdateIncomeRequest) error {
	ret := _m.Called(id, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, model.UpdateIncomeRequest) error); ok {
		r0 = rf(id, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
