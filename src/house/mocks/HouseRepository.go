// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	db "github.com/VlasovArtem/hob/src/db"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	model "github.com/VlasovArtem/hob/src/house/model"

	repository "github.com/VlasovArtem/hob/src/house/repository"

	uuid "github.com/google/uuid"
)

// HouseRepository is an autogenerated mock type for the HouseRepository type
type HouseRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: value, omit
func (_m *HouseRepository) Create(value interface{}, omit ...string) error {
	_va := make([]interface{}, len(omit))
	for _i := range omit {
		_va[_i] = omit[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, value)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, ...string) error); ok {
		r0 = rf(value, omit...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateBatch provides a mock function with given fields: entities
func (_m *HouseRepository) CreateBatch(entities []model.House) ([]model.House, error) {
	ret := _m.Called(entities)

	var r0 []model.House
	if rf, ok := ret.Get(0).(func([]model.House) []model.House); ok {
		r0 = rf(entities)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.House)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]model.House) error); ok {
		r1 = rf(entities)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DB provides a mock function with given fields:
func (_m *HouseRepository) DB() *gorm.DB {
	ret := _m.Called()

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func() *gorm.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// DBModeled provides a mock function with given fields: _a0
func (_m *HouseRepository) DBModeled(_a0 interface{}) *gorm.DB {
	ret := _m.Called(_a0)

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(interface{}) *gorm.DB); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// Delete provides a mock function with given fields: id
func (_m *HouseRepository) Delete(id uuid.UUID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields: id
func (_m *HouseRepository) Exists(id uuid.UUID) bool {
	ret := _m.Called(id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uuid.UUID) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ExistsBy provides a mock function with given fields: query, args
func (_m *HouseRepository) ExistsBy(query interface{}, args ...interface{}) bool {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) bool); ok {
		r0 = rf(query, args...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Find provides a mock function with given fields: id
func (_m *HouseRepository) Find(id uuid.UUID) (model.House, error) {
	ret := _m.Called(id)

	var r0 model.House
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.House); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.House)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindBy provides a mock function with given fields: query, conditions
func (_m *HouseRepository) FindBy(query interface{}, conditions ...interface{}) (model.House, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, conditions...)
	ret := _m.Called(_ca...)

	var r0 model.House
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) model.House); ok {
		r0 = rf(query, conditions...)
	} else {
		r0 = ret.Get(0).(model.House)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}, ...interface{}) error); ok {
		r1 = rf(query, conditions...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindById provides a mock function with given fields: id
func (_m *HouseRepository) FindById(id uuid.UUID) (model.House, error) {
	ret := _m.Called(id)

	var r0 model.House
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.House); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.House)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByUserId provides a mock function with given fields: id
func (_m *HouseRepository) FindByUserId(id uuid.UUID) []model.House {
	ret := _m.Called(id)

	var r0 []model.House
	if rf, ok := ret.Get(0).(func(uuid.UUID) []model.House); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.House)
		}
	}

	return r0
}

// FindHousesByGroupId provides a mock function with given fields: groupId
func (_m *HouseRepository) FindHousesByGroupId(groupId uuid.UUID) []model.House {
	ret := _m.Called(groupId)

	var r0 []model.House
	if rf, ok := ret.Get(0).(func(uuid.UUID) []model.House); ok {
		r0 = rf(groupId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.House)
		}
	}

	return r0
}

// FindHousesByGroupIds provides a mock function with given fields: groupIds
func (_m *HouseRepository) FindHousesByGroupIds(groupIds []uuid.UUID) []model.House {
	ret := _m.Called(groupIds)

	var r0 []model.House
	if rf, ok := ret.Get(0).(func([]uuid.UUID) []model.House); ok {
		r0 = rf(groupIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.House)
		}
	}

	return r0
}

// FindReceiver provides a mock function with given fields: receiver, id
func (_m *HouseRepository) FindReceiver(receiver interface{}, id uuid.UUID) error {
	ret := _m.Called(receiver, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, uuid.UUID) error); ok {
		r0 = rf(receiver, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindReceiverBy provides a mock function with given fields: receiver, query, conditions
func (_m *HouseRepository) FindReceiverBy(receiver interface{}, query interface{}, conditions ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, receiver, query)
	_ca = append(_ca, conditions...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, interface{}, ...interface{}) error); ok {
		r0 = rf(receiver, query, conditions...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// First provides a mock function with given fields: id
func (_m *HouseRepository) First(id uuid.UUID) (model.House, error) {
	ret := _m.Called(id)

	var r0 model.House
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.House); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.House)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FirstBy provides a mock function with given fields: query, conditions
func (_m *HouseRepository) FirstBy(query interface{}, conditions ...interface{}) (model.House, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, conditions...)
	ret := _m.Called(_ca...)

	var r0 model.House
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) model.House); ok {
		r0 = rf(query, conditions...)
	} else {
		r0 = ret.Get(0).(model.House)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}, ...interface{}) error); ok {
		r1 = rf(query, conditions...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FirstReceiver provides a mock function with given fields: receiver, id
func (_m *HouseRepository) FirstReceiver(receiver interface{}, id uuid.UUID) error {
	ret := _m.Called(receiver, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, uuid.UUID) error); ok {
		r0 = rf(receiver, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FirstReceiverBy provides a mock function with given fields: receiver, query, conditions
func (_m *HouseRepository) FirstReceiverBy(receiver interface{}, query interface{}, conditions ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, receiver, query)
	_ca = append(_ca, conditions...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, interface{}, ...interface{}) error); ok {
		r0 = rf(receiver, query, conditions...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetEntity provides a mock function with given fields:
func (_m *HouseRepository) GetEntity() model.House {
	ret := _m.Called()

	var r0 model.House
	if rf, ok := ret.Get(0).(func() model.House); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.House)
	}

	return r0
}

// GetProvider provides a mock function with given fields:
func (_m *HouseRepository) GetProvider() db.ProviderInterface {
	ret := _m.Called()

	var r0 db.ProviderInterface
	if rf, ok := ret.Get(0).(func() db.ProviderInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.ProviderInterface)
		}
	}

	return r0
}

// Modeled provides a mock function with given fields:
func (_m *HouseRepository) Modeled() *gorm.DB {
	ret := _m.Called()

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func() *gorm.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// Transactional provides a mock function with given fields: tx
func (_m *HouseRepository) Transactional(tx *gorm.DB) repository.HouseRepository {
	ret := _m.Called(tx)

	var r0 repository.HouseRepository
	if rf, ok := ret.Get(0).(func(*gorm.DB) repository.HouseRepository); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repository.HouseRepository)
		}
	}

	return r0
}

// Update provides a mock function with given fields: id, entity, omit
func (_m *HouseRepository) Update(id uuid.UUID, entity interface{}, omit ...string) error {
	_va := make([]interface{}, len(omit))
	for _i := range omit {
		_va[_i] = omit[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, id, entity)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, interface{}, ...string) error); ok {
		r0 = rf(id, entity, omit...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateByRequest provides a mock function with given fields: id, request
func (_m *HouseRepository) UpdateByRequest(id uuid.UUID, request model.UpdateHouseRequest) error {
	ret := _m.Called(id, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, model.UpdateHouseRequest) error); ok {
		r0 = rf(id, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewHouseRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewHouseRepository creates a new instance of HouseRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHouseRepository(t mockConstructorTestingTNewHouseRepository) *HouseRepository {
	mock := &HouseRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
