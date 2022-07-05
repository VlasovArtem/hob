// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	db "github.com/VlasovArtem/hob/src/db"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	model "github.com/VlasovArtem/hob/src/payment/model"

	repository "github.com/VlasovArtem/hob/src/payment/repository"

	time "time"

	uuid "github.com/google/uuid"
)

// PaymentRepository is an autogenerated mock type for the PaymentRepository type
type PaymentRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: value, omit
func (_m *PaymentRepository) Create(value interface{}, omit ...string) error {
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

// DB provides a mock function with given fields:
func (_m *PaymentRepository) DB() *gorm.DB {
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
func (_m *PaymentRepository) DBModeled(_a0 interface{}) *gorm.DB {
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
func (_m *PaymentRepository) Delete(id uuid.UUID) error {
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
func (_m *PaymentRepository) Exists(id uuid.UUID) bool {
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
func (_m *PaymentRepository) ExistsBy(query interface{}, args ...interface{}) bool {
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
func (_m *PaymentRepository) Find(id uuid.UUID) (model.Payment, error) {
	ret := _m.Called(id)

	var r0 model.Payment
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.Payment); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Payment)
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
func (_m *PaymentRepository) FindBy(query interface{}, conditions ...interface{}) (model.Payment, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, conditions...)
	ret := _m.Called(_ca...)

	var r0 model.Payment
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) model.Payment); ok {
		r0 = rf(query, conditions...)
	} else {
		r0 = ret.Get(0).(model.Payment)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}, ...interface{}) error); ok {
		r1 = rf(query, conditions...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByGroupId provides a mock function with given fields: id, limit, offset, from, to
func (_m *PaymentRepository) FindByGroupId(id uuid.UUID, limit int, offset int, from *time.Time, to *time.Time) []model.PaymentDto {
	ret := _m.Called(id, limit, offset, from, to)

	var r0 []model.PaymentDto
	if rf, ok := ret.Get(0).(func(uuid.UUID, int, int, *time.Time, *time.Time) []model.PaymentDto); ok {
		r0 = rf(id, limit, offset, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.PaymentDto)
		}
	}

	return r0
}

// FindByHouseId provides a mock function with given fields: houseId, limit, offset, from, to
func (_m *PaymentRepository) FindByHouseId(houseId uuid.UUID, limit int, offset int, from *time.Time, to *time.Time) []model.PaymentDto {
	ret := _m.Called(houseId, limit, offset, from, to)

	var r0 []model.PaymentDto
	if rf, ok := ret.Get(0).(func(uuid.UUID, int, int, *time.Time, *time.Time) []model.PaymentDto); ok {
		r0 = rf(houseId, limit, offset, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.PaymentDto)
		}
	}

	return r0
}

// FindByProviderId provides a mock function with given fields: providerId, limit, offset, from, to
func (_m *PaymentRepository) FindByProviderId(providerId uuid.UUID, limit int, offset int, from *time.Time, to *time.Time) []model.PaymentDto {
	ret := _m.Called(providerId, limit, offset, from, to)

	var r0 []model.PaymentDto
	if rf, ok := ret.Get(0).(func(uuid.UUID, int, int, *time.Time, *time.Time) []model.PaymentDto); ok {
		r0 = rf(providerId, limit, offset, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.PaymentDto)
		}
	}

	return r0
}

// FindByUserId provides a mock function with given fields: userId, limit, offset, from, to
func (_m *PaymentRepository) FindByUserId(userId uuid.UUID, limit int, offset int, from *time.Time, to *time.Time) []model.PaymentDto {
	ret := _m.Called(userId, limit, offset, from, to)

	var r0 []model.PaymentDto
	if rf, ok := ret.Get(0).(func(uuid.UUID, int, int, *time.Time, *time.Time) []model.PaymentDto); ok {
		r0 = rf(userId, limit, offset, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.PaymentDto)
		}
	}

	return r0
}

// FindReceiver provides a mock function with given fields: receiver, id
func (_m *PaymentRepository) FindReceiver(receiver interface{}, id uuid.UUID) error {
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
func (_m *PaymentRepository) FindReceiverBy(receiver interface{}, query interface{}, conditions ...interface{}) error {
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
func (_m *PaymentRepository) First(id uuid.UUID) (model.Payment, error) {
	ret := _m.Called(id)

	var r0 model.Payment
	if rf, ok := ret.Get(0).(func(uuid.UUID) model.Payment); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(model.Payment)
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
func (_m *PaymentRepository) FirstBy(query interface{}, conditions ...interface{}) (model.Payment, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, conditions...)
	ret := _m.Called(_ca...)

	var r0 model.Payment
	if rf, ok := ret.Get(0).(func(interface{}, ...interface{}) model.Payment); ok {
		r0 = rf(query, conditions...)
	} else {
		r0 = ret.Get(0).(model.Payment)
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
func (_m *PaymentRepository) FirstReceiver(receiver interface{}, id uuid.UUID) error {
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
func (_m *PaymentRepository) FirstReceiverBy(receiver interface{}, query interface{}, conditions ...interface{}) error {
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
func (_m *PaymentRepository) GetEntity() model.Payment {
	ret := _m.Called()

	var r0 model.Payment
	if rf, ok := ret.Get(0).(func() model.Payment); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.Payment)
	}

	return r0
}

// GetProvider provides a mock function with given fields:
func (_m *PaymentRepository) GetProvider() db.ProviderInterface {
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
func (_m *PaymentRepository) Modeled() *gorm.DB {
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
func (_m *PaymentRepository) Transactional(tx *gorm.DB) repository.PaymentRepository {
	ret := _m.Called(tx)

	var r0 repository.PaymentRepository
	if rf, ok := ret.Get(0).(func(*gorm.DB) repository.PaymentRepository); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repository.PaymentRepository)
		}
	}

	return r0
}

// Update provides a mock function with given fields: id, entity, omit
func (_m *PaymentRepository) Update(id uuid.UUID, entity interface{}, omit ...string) error {
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

type mockConstructorTestingTNewPaymentRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewPaymentRepository creates a new instance of PaymentRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPaymentRepository(t mockConstructorTestingTNewPaymentRepository) *PaymentRepository {
	mock := &PaymentRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
