package mock

import (
	"context"
	"github.com/google/uuid"
	"github.com/robfig/cron"
	"github.com/stretchr/testify/mock"
	hm "house/model"
	i "income/model"
	isr "income/scheduler/model"
	m "meter/model"
	pm "payment/model"
	psr "payment/scheduler/model"
	um "user/model"
)

type UserServiceMock struct {
	mock.Mock
}

func (u *UserServiceMock) Add(request um.CreateUserRequest) (um.UserResponse, error) {
	args := u.Called(request)

	return args.Get(0).(um.UserResponse), args.Error(1)
}

func (u *UserServiceMock) FindById(id uuid.UUID) (um.UserResponse, error) {
	args := u.Called(id)

	return args.Get(0).(um.UserResponse), args.Error(1)
}

func (u *UserServiceMock) ExistsById(id uuid.UUID) bool {
	args := u.Called(id)

	return args.Bool(0)
}

func (u *UserServiceMock) ExistsByEmail(email string) bool {
	args := u.Called(email)

	return args.Bool(0)
}

type HouseServiceMock struct {
	mock.Mock
}

func (h *HouseServiceMock) Add(house hm.CreateHouseRequest) (hm.HouseResponse, error) {
	args := h.Called(house)

	return args.Get(0).(hm.HouseResponse), args.Error(0)
}

func (h *HouseServiceMock) FindAll() []hm.HouseResponse {
	args := h.Called()

	return args.Get(0).([]hm.HouseResponse)
}

func (h *HouseServiceMock) FindById(id uuid.UUID) (hm.HouseResponse, error) {
	args := h.Called(id)

	return args.Get(0).(hm.HouseResponse), args.Error(0)
}

func (h *HouseServiceMock) ExistsById(id uuid.UUID) bool {
	args := h.Called(id)

	return args.Bool(0)
}

type PaymentServiceMock struct {
	mock.Mock
}

func (p *PaymentServiceMock) Add(request pm.CreatePaymentRequest) (pm.PaymentResponse, error) {
	args := p.Called(request)

	return args.Get(0).(pm.PaymentResponse), args.Error(1)
}

func (p *PaymentServiceMock) FindById(id uuid.UUID) (pm.PaymentResponse, error) {
	args := p.Called(id)

	return args.Get(0).(pm.PaymentResponse), args.Error(1)
}

func (p *PaymentServiceMock) FindByHouseId(houseId uuid.UUID) []pm.PaymentResponse {
	args := p.Called(houseId)

	return args.Get(0).([]pm.PaymentResponse)
}

func (p *PaymentServiceMock) FindByUserId(userId uuid.UUID) []pm.PaymentResponse {
	args := p.Called(userId)

	return args.Get(0).([]pm.PaymentResponse)
}

func (p *PaymentServiceMock) ExistsById(id uuid.UUID) bool {
	args := p.Called(id)

	return args.Bool(0)
}

type MeterServiceMock struct {
	mock.Mock
}

func (ms *MeterServiceMock) AddMeter(request m.CreateMeterRequest) (m.MeterResponse, error) {
	args := ms.Called(request)

	return args.Get(0).(m.MeterResponse), args.Error(1)
}

func (ms *MeterServiceMock) FindById(id uuid.UUID) (m.MeterResponse, error) {
	args := ms.Called(id)

	return args.Get(0).(m.MeterResponse), args.Error(1)
}

func (ms *MeterServiceMock) FindByPaymentId(id uuid.UUID) (m.MeterResponse, error) {
	args := ms.Called(id)

	return args.Get(0).(m.MeterResponse), args.Error(1)
}

type IncomeServiceMock struct {
	mock.Mock
}

func (is *IncomeServiceMock) Add(request i.CreateIncomeRequest) (i.IncomeResponse, error) {
	args := is.Called(request)

	return args.Get(0).(i.IncomeResponse), args.Error(1)
}

func (is *IncomeServiceMock) FindById(id uuid.UUID) (i.IncomeResponse, error) {
	args := is.Called(id)

	return args.Get(0).(i.IncomeResponse), args.Error(1)
}

func (is *IncomeServiceMock) FindByHouseId(id uuid.UUID) (i.IncomeResponse, error) {
	args := is.Called(id)

	return args.Get(0).(i.IncomeResponse), args.Error(1)
}

type SchedulerServiceMock struct {
	mock.Mock
}

func (s *SchedulerServiceMock) Add(scheduledItemId uuid.UUID, scheduleSpec string, scheduleFunc func()) (cron.EntryID, error) {
	args := s.Called(scheduledItemId, scheduleSpec, scheduleFunc)

	return args.Get(0).(cron.EntryID), args.Error(1)
}

func (s *SchedulerServiceMock) Remove(id uuid.UUID) error {
	args := s.Called(id)

	return args.Error(0)
}

func (s *SchedulerServiceMock) Stop() context.Context {
	args := s.Called()

	return args.Get(0).(context.Context)
}

type PaymentSchedulerServiceMock struct {
	mock.Mock
}

func (p *PaymentSchedulerServiceMock) Add(request psr.CreatePaymentSchedulerRequest) (psr.PaymentSchedulerResponse, error) {
	args := p.Called(request)

	return args.Get(0).(psr.PaymentSchedulerResponse), args.Error(1)
}

func (p *PaymentSchedulerServiceMock) Remove(id uuid.UUID) error {
	args := p.Called(id)

	return args.Error(0)
}

func (p *PaymentSchedulerServiceMock) FindById(id uuid.UUID) (psr.PaymentSchedulerResponse, error) {
	args := p.Called(id)

	return args.Get(0).(psr.PaymentSchedulerResponse), args.Error(1)
}

func (p *PaymentSchedulerServiceMock) FindByHouseId(houseId uuid.UUID) []psr.PaymentSchedulerResponse {
	args := p.Called(houseId)

	return args.Get(0).([]psr.PaymentSchedulerResponse)
}

func (p *PaymentSchedulerServiceMock) FindByUserId(userId uuid.UUID) []psr.PaymentSchedulerResponse {
	args := p.Called(userId)

	return args.Get(0).([]psr.PaymentSchedulerResponse)
}

type IncomeSchedulerServiceMock struct {
	mock.Mock
}

func (i *IncomeSchedulerServiceMock) Add(request isr.CreateIncomeSchedulerRequest) (isr.IncomeSchedulerResponse, error) {
	args := i.Called(request)

	return args.Get(0).(isr.IncomeSchedulerResponse), args.Error(1)
}

func (i *IncomeSchedulerServiceMock) Remove(id uuid.UUID) error {
	args := i.Called(id)

	return args.Error(0)
}

func (i *IncomeSchedulerServiceMock) FindById(id uuid.UUID) (isr.IncomeSchedulerResponse, error) {
	args := i.Called(id)

	return args.Get(0).(isr.IncomeSchedulerResponse), args.Error(1)
}

func (i *IncomeSchedulerServiceMock) FindByHouseId(houseId uuid.UUID) (isr.IncomeSchedulerResponse, error) {
	args := i.Called(houseId)

	return args.Get(0).(isr.IncomeSchedulerResponse), args.Error(1)
}
