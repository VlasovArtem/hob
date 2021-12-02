// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	model "github.com/VlasovArtem/hob/src/house/model"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// HouseRepository is an autogenerated mock type for the HouseRepository type
type HouseRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: entity
func (_m *HouseRepository) Create(entity model.House) (model.House, error) {
	ret := _m.Called(entity)

	var r0 model.House
	if rf, ok := ret.Get(0).(func(model.House) model.House); ok {
		r0 = rf(entity)
	} else {
		r0 = ret.Get(0).(model.House)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.House) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExistsById provides a mock function with given fields: id
func (_m *HouseRepository) ExistsById(id uuid.UUID) bool {
	ret := _m.Called(id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uuid.UUID) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FindResponseById provides a mock function with given fields: id
func (_m *HouseRepository) FindResponseById(id uuid.UUID) (model.HouseDto, error) {
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

// FindResponseByUserId provides a mock function with given fields: id
func (_m *HouseRepository) FindResponseByUserId(id uuid.UUID) []model.HouseDto {
	ret := _m.Called(id)

	var r0 []model.HouseDto
	if rf, ok := ret.Get(0).(func(uuid.UUID) []model.HouseDto); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.HouseDto)
		}
	}

	return r0
}