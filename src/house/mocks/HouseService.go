// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	model "house/model"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// HouseService is an autogenerated mock type for the HouseService type
type HouseService struct {
	mock.Mock
}

// Add provides a mock function with given fields: house
func (_m *HouseService) Add(house model.CreateHouseRequest) (model.HouseDto, error) {
	ret := _m.Called(house)

	var r0 model.HouseDto
	if rf, ok := ret.Get(0).(func(model.CreateHouseRequest) model.HouseDto); ok {
		r0 = rf(house)
	} else {
		r0 = ret.Get(0).(model.HouseDto)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.CreateHouseRequest) error); ok {
		r1 = rf(house)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExistsById provides a mock function with given fields: id
func (_m *HouseService) ExistsById(id uuid.UUID) bool {
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
func (_m *HouseService) FindById(id uuid.UUID) (model.HouseDto, error) {
	ret := _m.Called(id)

	var r0 model.HouseDto
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.HouseDto); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.HouseDto)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByUserId provides a mock function with given fields: userId
func (_m *HouseService) FindByUserId(userId uuid.UUID) []model.HouseDto {
	ret := _m.Called(userId)

	var r0 []model.HouseDto
	if rf, ok := ret.Get(0).(func(uuid.UUID) []model.HouseDto); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.HouseDto)
		}
	}

	return r0
}
