package repository

import (
	"github.com/VlasovArtem/hob/src/db"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/payment/mocks"
	"github.com/VlasovArtem/hob/src/payment/model"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type PaymentRepositoryTestSuite struct {
	database.DBTestSuite
	repository   PaymentRepository
	createdUser  userModel.User
	createdHouse houseModel.House
}

func (p *PaymentRepositoryTestSuite) SetupSuite() {
	p.InitDBTestSuite()

	p.CreateRepository(
		func(service db.DatabaseService) {
			p.repository = NewPaymentRepository(service)
		},
	).
		AddMigrators(userModel.User{}, houseModel.House{}, model.Payment{})

	p.createdUser = userMocks.GenerateUser()
	p.CreateEntity(&p.createdUser)

	p.createdHouse = houseMocks.GenerateHouse(p.createdUser.Id)
	p.CreateEntity(&p.createdHouse)
}

func (p *PaymentRepositoryTestSuite) TearDownSuite() {
	p.TearDown()
}

func TestPaymentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentRepositoryTestSuite))
}

func (p *PaymentRepositoryTestSuite) Test_Create() {
	payment := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id)

	actual, err := p.repository.Create(payment)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), payment, actual)
}

func (p *PaymentRepositoryTestSuite) Test_Creat_WithMissingUser() {
	payment := mocks.GeneratePayment(p.createdHouse.Id, uuid.New())

	actual, err := p.repository.Create(payment)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), payment, actual)
}

func (p *PaymentRepositoryTestSuite) Test_Creat_WithMissingHouse() {
	payment := mocks.GeneratePayment(uuid.New(), p.createdUser.Id)

	actual, err := p.repository.Create(payment)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), payment, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindById() {
	payment := p.createPayment()

	actual, err := p.repository.FindById(payment.Id)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), payment, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindById_WithMissingId() {
	actual, err := p.repository.FindById(uuid.New())

	assert.ErrorIs(p.T(), err, gorm.ErrRecordNotFound)
	assert.Equal(p.T(), model.Payment{}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByUserId() {
	payment := p.createPayment()

	actual := p.repository.FindByUserId(payment.UserId)

	var actualResponse model.Payment

	for _, response := range actual {
		if response.Id == payment.Id {
			actualResponse = response
			break
		}
	}
	assert.Equal(p.T(), payment, actualResponse)
}

func (p *PaymentRepositoryTestSuite) Test_FindByUserId_WithMissingUserId() {
	actual := p.repository.FindByUserId(uuid.New())

	assert.Equal(p.T(), []model.Payment{}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByHouseId() {
	payment := p.createPayment()

	actual := p.repository.FindByHouseId(payment.HouseId)

	var actualResponse model.Payment

	for _, response := range actual {
		if response.Id == payment.Id {
			actualResponse = response
			break
		}
	}
	assert.Equal(p.T(), payment, actualResponse)
}

func (p *PaymentRepositoryTestSuite) Test_FindByHouseId_WithMissingUserId() {
	actual := p.repository.FindByHouseId(uuid.New())

	assert.Equal(p.T(), []model.Payment{}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_ExistsById() {
	payment := p.createPayment()

	assert.True(p.T(), p.repository.ExistsById(payment.Id))
}

func (p *PaymentRepositoryTestSuite) Test_ExistsById_WithMissingId() {
	assert.False(p.T(), p.repository.ExistsById(uuid.New()))
}

func (p *PaymentRepositoryTestSuite) createPayment() model.Payment {
	payment := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id)

	create, err := p.repository.Create(payment)

	assert.Nil(p.T(), err)

	return create
}
