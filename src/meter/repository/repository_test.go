package repository

import (
	"github.com/VlasovArtem/hob/src/db"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	meterMocks "github.com/VlasovArtem/hob/src/meter/mocks"
	"github.com/VlasovArtem/hob/src/meter/model"
	paymentMocks "github.com/VlasovArtem/hob/src/payment/mocks"
	paymentModel "github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type MeterRepositoryTestSuite struct {
	database.DBTestSuite
	repository     MeterRepository
	createdPayment paymentModel.Payment
	createdUser    userModel.User
	createdHouse   houseModel.House
}

func (m *MeterRepositoryTestSuite) SetupSuite() {
	m.InitDBTestSuite()

	m.CreateRepository(
		func(service db.DatabaseService) {
			m.repository = NewMeterRepository(service)
		},
	).
		AddMigrators(userModel.User{}, houseModel.House{}, paymentModel.Payment{}, model.Meter{})

	m.createdUser = userMocks.GenerateUser()
	m.CreateEntity(&m.createdUser)

	m.AddBeforeTest(
		func(service db.DatabaseService) {
			m.createdHouse = houseMocks.GenerateHouse(m.createdUser.Id)
			m.CreateEntity(&m.createdHouse)
		})
	m.AddBeforeTest(
		func(service db.DatabaseService) {
			m.createdPayment = paymentMocks.GeneratePayment(m.createdHouse.Id, m.createdUser.Id)
			m.CreateEntity(&m.createdPayment)
		})
}

func (m *MeterRepositoryTestSuite) TearDownSuite() {
	m.TearDown()
}

func TestPaymentRepositorySchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(MeterRepositoryTestSuite))
}

func (m *MeterRepositoryTestSuite) Test_Create() {
	meter := meterMocks.GenerateMeter(m.createdPayment.Id, m.createdHouse.Id)

	actual, err := m.repository.Create(meter)

	assert.Nil(m.T(), err)
	assert.Equal(m.T(), meter, actual)
}

func (m *MeterRepositoryTestSuite) Test_Creat_WithMissingPayment() {
	meter := meterMocks.GenerateMeter(uuid.New(), m.createdHouse.Id)

	actual, err := m.repository.Create(meter)

	assert.NotNil(m.T(), err)
	assert.Equal(m.T(), meter, actual)
}

func (m *MeterRepositoryTestSuite) Test_Creat_WithMissingHouse() {
	meter := meterMocks.GenerateMeter(m.createdUser.Id, uuid.New())

	actual, err := m.repository.Create(meter)

	assert.NotNil(m.T(), err)
	assert.Equal(m.T(), meter, actual)
}

func (m *MeterRepositoryTestSuite) Test_FindById() {
	meter := m.createMeter()

	actual, err := m.repository.FindById(meter.Id)

	assert.Nil(m.T(), err)
	assert.Equal(m.T(), meter, actual)
}

func (m *MeterRepositoryTestSuite) Test_FindById_WithMissingId() {
	actual, err := m.repository.FindById(uuid.New())

	assert.ErrorIs(m.T(), err, gorm.ErrRecordNotFound)
	assert.Equal(m.T(), model.Meter{}, actual)
}

func (m *MeterRepositoryTestSuite) Test_FindByPaymentId() {
	meter := m.createMeter()

	meterResponse, err := m.repository.FindByPaymentId(meter.PaymentId)

	assert.Nil(m.T(), err)
	assert.Equal(m.T(), meter, meterResponse)
}

func (m *MeterRepositoryTestSuite) Test_FindByPaymentId_WithMissingUserId() {
	meterResponse, err := m.repository.FindByPaymentId(uuid.New())

	assert.Equal(m.T(), gorm.ErrRecordNotFound, err)
	assert.Equal(m.T(), model.Meter{}, meterResponse)
}

func (m *MeterRepositoryTestSuite) Test_ExistsById() {
	payment := m.createMeter()

	assert.True(m.T(), m.repository.ExistsById(payment.Id))
}

func (m *MeterRepositoryTestSuite) Test_ExistsById_WithMissingId() {
	assert.False(m.T(), m.repository.ExistsById(uuid.New()))
}

func (m *MeterRepositoryTestSuite) Test_FindByHouseId() {
	meter := m.createMeter()

	meters := m.repository.FindByHouseId(meter.HouseId)

	assert.Equal(m.T(), []model.Meter{meter}, meters)
}

func (m *MeterRepositoryTestSuite) Test_FindByHouseId_WithMissingRecords() {
	meters := m.repository.FindByHouseId(uuid.New())

	assert.Equal(m.T(), []model.Meter{}, meters)
}

func (m *MeterRepositoryTestSuite) createMeter() model.Meter {
	meter := meterMocks.GenerateMeter(m.createdPayment.Id, m.createdHouse.Id)

	create, err := m.repository.Create(meter)

	assert.Nil(m.T(), err)

	return create
}
