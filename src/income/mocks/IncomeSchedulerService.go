// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	model "github.com/VlasovArtem/hob/src/income/scheduler/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// IncomeSchedulerService is an autogenerated mock type for the IncomeSchedulerService type
type IncomeSchedulerService struct {
	mock.Mock
}

// Add provides a mock function with given fields: request
func (_m *IncomeSchedulerService) Add(request model.CreateIncomeSchedulerRequest) (model.IncomeSchedulerDto, error) {
	ret := _m.Called(request)

	var r0 model.IncomeSchedulerDto
	if rf, ok := ret.Get(0).(func(model.CreateIncomeSchedulerRequest) model.IncomeSchedulerDto); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Get(0).(model.IncomeSchedulerDto)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.CreateIncomeSchedulerRequest) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteById provides a mock function with given fields: id
func (_m *IncomeSchedulerService) DeleteById(id uuid.UUID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByHouseId provides a mock function with given fields: id
func (_m *IncomeSchedulerService) FindByHouseId(id uuid.UUID) []model.IncomeSchedulerDto {
	ret := _m.Called(id)

	var r0 []model.IncomeSchedulerDto
	if rf, ok := ret.Get(0).(func(uuid.UUID) []model.IncomeSchedulerDto); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.IncomeSchedulerDto)
		}
	}

	return r0
}

// FindById provides a mock function with given fields: id
func (_m *IncomeSchedulerService) FindById(id uuid.UUID) (model.IncomeSchedulerDto, error) {
	ret := _m.Called(id)

	var r0 model.IncomeSchedulerDto
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.IncomeSchedulerDto); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.IncomeSchedulerDto)
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
func (_m *IncomeSchedulerService) Update(id uuid.UUID, request model.UpdateIncomeSchedulerRequest) error {
	ret := _m.Called(id, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, model.UpdateIncomeSchedulerRequest) error); ok {
		r0 = rf(id, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
