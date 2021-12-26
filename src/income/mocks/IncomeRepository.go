// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	model "github.com/VlasovArtem/hob/src/income/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// IncomeRepository is an autogenerated mock type for the IncomeRepository type
type IncomeRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: entity
func (_m *IncomeRepository) Create(entity model.Income) (model.Income, error) {
	ret := _m.Called(entity)

	var r0 model.Income
	if rf, ok := ret.Get(0).(func(model.Income) model.Income); ok {
		r0 = rf(entity)
	} else {
		r0 = ret.Get(0).(model.Income)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.Income) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteById provides a mock function with given fields: id
func (_m *IncomeRepository) DeleteById(id uuid.UUID) error {
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
func (_m *IncomeRepository) ExistsById(id uuid.UUID) bool {
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
func (_m *IncomeRepository) FindByHouseId(id uuid.UUID) ([]model.IncomeDto, error) {
	ret := _m.Called(id)

	var r0 []model.IncomeDto
	if rf, ok := ret.Get(0).(func(uuid.UUID) []model.IncomeDto); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.IncomeDto)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindById provides a mock function with given fields: id
func (_m *IncomeRepository) FindById(id uuid.UUID) (model.Income, error) {
	ret := _m.Called(id)

	var r0 model.Income
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.Income); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Income)
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
func (_m *IncomeRepository) Update(id uuid.UUID, request model.UpdateIncomeRequest) error {
	ret := _m.Called(id, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, model.UpdateIncomeRequest) error); ok {
		r0 = rf(id, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
