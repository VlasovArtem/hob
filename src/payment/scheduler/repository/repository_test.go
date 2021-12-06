package repository

import (
	"github.com/VlasovArtem/hob/src/db"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/payment/scheduler/mocks"
	"github.com/VlasovArtem/hob/src/payment/scheduler/model"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type PaymentRepositorySchedulerTestSuite struct {
	database.DBTestSuite
	repository   PaymentSchedulerRepository
	createdUser  userModel.User
	createdHouse houseModel.House
}

func (p *PaymentRepositorySchedulerTestSuite) SetupSuite() {
	p.InitDBTestSuite()

	p.CreateRepository(
		func(service db.DatabaseService) {
			p.repository = NewPaymentSchedulerRepository(service)
		},
	).
		AddMigrators(userModel.User{}, houseModel.House{}, model.PaymentScheduler{})

	p.createdUser = userMocks.GenerateUser()
	p.CreateEntity(&p.createdUser)

	p.createdHouse = houseMocks.GenerateHouse(p.createdUser.Id)
	p.CreateEntity(&p.createdHouse)
}

func (p *PaymentRepositorySchedulerTestSuite) TearDownSuite() {
	p.TearDown()
}

func TestPaymentRepositorySchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentRepositorySchedulerTestSuite))
}

func (p *PaymentRepositorySchedulerTestSuite) Test_Create() {
	paymentScheduler := mocks.GeneratePaymentScheduler(p.createdHouse.Id, p.createdUser.Id)

	actual, err := p.repository.Create(paymentScheduler)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), paymentScheduler, actual)

	p.Delete(paymentScheduler)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_Creat_WithMissingUser() {
	paymentScheduler := mocks.GeneratePaymentScheduler(p.createdHouse.Id, uuid.New())

	actual, err := p.repository.Create(paymentScheduler)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), paymentScheduler, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_Creat_WithMissingHouse() {
	paymentScheduler := mocks.GeneratePaymentScheduler(uuid.New(), p.createdUser.Id)

	actual, err := p.repository.Create(paymentScheduler)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), paymentScheduler, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_FindById() {
	payment := p.createPaymentScheduler()

	actual, err := p.repository.FindById(payment.Id)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), payment, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_FindById_WithMissingId() {
	actual, err := p.repository.FindById(uuid.New())

	assert.ErrorIs(p.T(), err, gorm.ErrRecordNotFound)
	assert.Equal(p.T(), model.PaymentScheduler{}, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_FindByUserId() {
	payment := p.createPaymentScheduler()

	actual := p.repository.FindByUserId(payment.UserId)

	var actualResponse model.PaymentScheduler

	for _, response := range actual {
		if response.Id == payment.Id {
			actualResponse = response
			break
		}
	}
	assert.Equal(p.T(), payment, actualResponse)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_FindByUserId_WithMissingId() {
	actual := p.repository.FindByUserId(uuid.New())

	assert.Equal(p.T(), []model.PaymentScheduler{}, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_FindByHouseId() {
	payment := p.createPaymentScheduler()

	actual := p.repository.FindByHouseId(payment.HouseId)

	var actualResponse model.PaymentScheduler

	for _, response := range actual {
		if response.Id == payment.Id {
			actualResponse = response
			break
		}
	}
	assert.Equal(p.T(), payment, actualResponse)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_FindByHouseId_WithMissingId() {
	actual := p.repository.FindByHouseId(uuid.New())

	assert.Equal(p.T(), []model.PaymentScheduler{}, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_ExistsById() {
	payment := p.createPaymentScheduler()

	assert.True(p.T(), p.repository.ExistsById(payment.Id))
}

func (p *PaymentRepositorySchedulerTestSuite) Test_ExistsById_WithMissingId() {
	assert.False(p.T(), p.repository.ExistsById(uuid.New()))
}

func (p *PaymentRepositorySchedulerTestSuite) Test_DeleteById() {
	payment := p.createPaymentScheduler()

	assert.True(p.T(), p.repository.ExistsById(payment.Id))

	p.repository.DeleteById(payment.Id)

	assert.False(p.T(), p.repository.ExistsById(payment.Id))
}

func (p *PaymentRepositorySchedulerTestSuite) Test_DeleteById_WithMissingId() {
	p.repository.DeleteById(uuid.New())
}

func (p *PaymentRepositorySchedulerTestSuite) createPaymentScheduler() model.PaymentScheduler {
	payment := mocks.GeneratePaymentScheduler(p.createdHouse.Id, p.createdUser.Id)

	p.CreateEntity(payment)

	return payment
}
