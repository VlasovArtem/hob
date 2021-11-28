package repository

import (
	"db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	houseMocks "house/mocks"
	houseModel "house/model"
	"log"
	meterMocks "meter/mocks"
	"meter/model"
	paymentMocks "payment/mocks"
	paymentModel "payment/model"
	"test/testhelper/database"
	"testing"
	userMocks "user/mocks"
	userModel "user/model"
)

type MeterRepositoryTestSuite struct {
	suite.Suite
	database       db.DatabaseService
	repository     MeterRepository
	createdPayment paymentModel.Payment
	createdUser    userModel.User
	createdHouse   houseModel.House
}

func (m *MeterRepositoryTestSuite) SetupSuite() {
	config := db.NewDefaultDatabaseConfiguration()
	config.DBName = "hob_test"
	m.database = db.NewDatabaseService(config)
	m.repository = NewMeterRepository(m.database)
	err := m.database.D().AutoMigrate(model.Meter{})

	if err != nil {
		log.Fatal(err)
	}

	m.createdUser = userMocks.GenerateUser()
	err = m.database.Create(&m.createdUser)

	if err != nil {
		log.Fatal(err)
	}
}

func (m *MeterRepositoryTestSuite) TearDownSuite() {
	database.DropTable(m.database.D(), houseModel.House{})
	database.DropTable(m.database.D(), userModel.User{})
	database.DropTable(m.database.D(), paymentModel.Payment{})
	database.DropTable(m.database.D(), model.Meter{})
}

func (m *MeterRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	m.createdHouse = houseMocks.GenerateHouse(m.createdUser.Id)

	err := m.database.Create(&m.createdHouse)

	if err != nil {
		log.Fatal(err)
	}

	m.createdPayment = paymentMocks.GeneratePayment(m.createdHouse.Id, m.createdUser.Id)

	err = m.database.Create(&m.createdPayment)

	if err != nil {
		log.Fatal(err)
	}
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
