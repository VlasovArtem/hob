// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	cron "github.com/robfig/cron/v3"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// ServiceScheduler is an autogenerated mock type for the ServiceScheduler type
type ServiceScheduler struct {
	mock.Mock
}

// Add provides a mock function with given fields: scheduledItemId, scheduleSpec, scheduleFunc
func (_m *ServiceScheduler) Add(scheduledItemId uuid.UUID, scheduleSpec string, scheduleFunc func()) (cron.EntryID, error) {
	ret := _m.Called(scheduledItemId, scheduleSpec, scheduleFunc)

	var r0 cron.EntryID
	if rf, ok := ret.Get(0).(func(uuid.UUID, string, func()) cron.EntryID); ok {
		r0 = rf(scheduledItemId, scheduleSpec, scheduleFunc)
	} else {
		r0 = ret.Get(0).(cron.EntryID)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string, func()) error); ok {
		r1 = rf(scheduledItemId, scheduleSpec, scheduleFunc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Remove provides a mock function with given fields: id
func (_m *ServiceScheduler) Remove(id uuid.UUID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *ServiceScheduler) Stop() context.Context {
	ret := _m.Called()

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// Update provides a mock function with given fields: id, scheduleSpec, scheduleFunc
func (_m *ServiceScheduler) Update(id uuid.UUID, scheduleSpec string, scheduleFunc func()) (cron.EntryID, error) {
	ret := _m.Called(id, scheduleSpec, scheduleFunc)

	var r0 cron.EntryID
	if rf, ok := ret.Get(0).(func(uuid.UUID, string, func()) cron.EntryID); ok {
		r0 = rf(id, scheduleSpec, scheduleFunc)
	} else {
		r0 = ret.Get(0).(cron.EntryID)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, string, func()) error); ok {
		r1 = rf(id, scheduleSpec, scheduleFunc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}