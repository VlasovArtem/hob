package mock

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	hm "house/model"
	m "meter/model"
	pm "payment/model"
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

func (p *PaymentServiceMock) AddPayment(request pm.CreatePaymentRequest) (pm.PaymentResponse, error) {
	args := p.Called(request)

	return args.Get(0).(pm.PaymentResponse), args.Error(1)
}

func (p *PaymentServiceMock) FindPaymentById(id uuid.UUID) (pm.PaymentResponse, error) {
	args := p.Called(id)

	return args.Get(0).(pm.PaymentResponse), args.Error(1)
}

func (p *PaymentServiceMock) FindPaymentByHouseId(houseId uuid.UUID) []pm.PaymentResponse {
	args := p.Called(houseId)

	return args.Get(0).([]pm.PaymentResponse)
}

func (p *PaymentServiceMock) FindPaymentByUserId(userId uuid.UUID) []pm.PaymentResponse {
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
