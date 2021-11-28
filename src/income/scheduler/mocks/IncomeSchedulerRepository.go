// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	model "income/scheduler/model"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// IncomeSchedulerRepository is an autogenerated mock type for the IncomeSchedulerRepository type
type IncomeSchedulerRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: scheduler
func (_m *IncomeSchedulerRepository) Create(scheduler model.IncomeScheduler) (model.IncomeScheduler, error) {
	ret := _m.Called(scheduler)

	var r0 model.IncomeScheduler
	if rf, ok := ret.Get(0).(func(model.IncomeScheduler) model.IncomeScheduler); ok {
		r0 = rf(scheduler)
	} else {
		r0 = ret.Get(0).(model.IncomeScheduler)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.IncomeScheduler) error); ok {
		r1 = rf(scheduler)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteById provides a mock function with given fields: id
func (_m *IncomeSchedulerRepository) DeleteById(id uuid.UUID) {
	_m.Called(id)
}

// ExistsById provides a mock function with given fields: id
func (_m *IncomeSchedulerRepository) ExistsById(id uuid.UUID) bool {
	ret := _m.Called(id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uuid.UUID) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FindByHouseId provides a mock function with given fields: houseId
func (_m *IncomeSchedulerRepository) FindByHouseId(houseId uuid.UUID) []model.IncomeScheduler {
	ret := _m.Called(houseId)

	var r0 []model.IncomeScheduler
	if rf, ok := ret.Get(0).(func(uuid.UUID) []model.IncomeScheduler); ok {
		r0 = rf(houseId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.IncomeScheduler)
		}
	}

	return r0
}

// FindById provides a mock function with given fields: id
func (_m *IncomeSchedulerRepository) FindById(id uuid.UUID) (model.IncomeScheduler, error) {
	ret := _m.Called(id)

	var r0 model.IncomeScheduler
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.IncomeScheduler); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.IncomeScheduler)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
