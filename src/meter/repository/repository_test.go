package repository

import (
	"encoding/json"
	"github.com/VlasovArtem/hob/src/db"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	meterMocks "github.com/VlasovArtem/hob/src/meter/mocks"
	"github.com/VlasovArtem/hob/src/meter/model"
	paymentMocks "github.com/VlasovArtem/hob/src/payment/mocks"
	paymentModel "github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/provider/mocks"
	providerModel "github.com/VlasovArtem/hob/src/provider/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
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
	database.DBTestSuite[model.Meter]
	repository      MeterRepository
	createdProvider providerModel.Provider
	createdPayment  paymentModel.Payment
	createdUser     userModel.User
	createdHouse    houseModel.House
}

func (m *MeterRepositoryTestSuite) SetupSuite() {
	m.InitDBTestSuite()

	m.CreateRepository(
		func(service db.ModeledDatabase[model.Meter]) {
			m.repository = NewMeterRepository(service)
		},
	).
		AddAfterTest(func(service db.ModeledDatabase[model.Meter]) {
			testhelper.TruncateTable(service, model.Meter{})
			testhelper.TruncateTable(service, paymentModel.Payment{})
		}).
		AddAfterSuite(func(service db.ModeledDatabase[model.Meter]) {
			testhelper.TruncateTable(service, providerModel.Provider{})
			testhelper.TruncateTable(service, houseModel.House{})
			testhelper.TruncateTable(service, userModel.User{})
		}).
		ExecuteMigration(userModel.User{}, houseModel.House{}, providerModel.Provider{}, paymentModel.Payment{}, model.Meter{})

	m.createdUser = userMocks.GenerateUser()
	m.CreateEntity(&m.createdUser)

	m.createdProvider = mocks.GenerateProvider(m.createdUser.Id)
	m.CreateEntity(&m.createdProvider)

	m.createdHouse = houseMocks.GenerateHouse(m.createdUser.Id)
	m.CreateEntity(&m.createdHouse)

	m.AddBeforeTest(
		func(service db.ModeledDatabase[model.Meter]) {
			m.createdPayment = paymentMocks.GeneratePayment(m.createdHouse.Id, m.createdUser.Id, m.createdProvider.Id)
			m.CreateEntity(&m.createdPayment)
		})
}

func TestPaymentRepositorySchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(MeterRepositoryTestSuite))
}

func (m *MeterRepositoryTestSuite) Test_FindByPaymentId() {
	meter := m.createMeter()

	meterResponse, err := m.repository.FindByPaymentId(meter.PaymentId)

	assert.Nil(m.T(), err)
	assert.Equal(m.T(), meter, meterResponse)
}

func (m *MeterRepositoryTestSuite) Test_FindByPaymentId_WithMissingId() {
	meterResponse, err := m.repository.FindByPaymentId(uuid.New())

	assert.Equal(m.T(), gorm.ErrRecordNotFound, err)
	assert.Equal(m.T(), model.Meter{}, meterResponse)
}

func (m *MeterRepositoryTestSuite) Test_Update() {
	meter := m.createMeter()

	newDetails := map[string]float64{
		"first":  1.1,
		"second": 2.2,
		"third":  3.0,
	}

	marshal, _ := json.Marshal(newDetails)

	err := m.repository.Update(meter.Id, model.Meter{
		Name:        "Name New",
		Description: "Details New",
		Details:     marshal,
	})

	assert.Nil(m.T(), err)

	updatedMeter, err := m.repository.Find(meter.Id)

	assert.Nil(m.T(), err)

	assert.Equal(m.T(), model.Meter{
		Id:          meter.Id,
		Name:        "Name New",
		Description: "Details New",
		Details:     marshal,
		PaymentId:   meter.PaymentId,
		Payment:     meter.Payment,
	}, updatedMeter)
}

func (m *MeterRepositoryTestSuite) createMeter() model.Meter {
	meter := meterMocks.GenerateMeter(m.createdPayment.Id)

	m.CreateEntity(meter)

	return meter
}
