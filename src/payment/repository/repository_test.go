package repository

import (
	"fmt"
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

	p.Delete(payment)
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

func (p *PaymentRepositoryTestSuite) Test_DeleteById() {
	payment := p.createPayment()

	assert.Nil(p.T(), p.repository.DeleteById(payment.Id))
}

func (p *PaymentRepositoryTestSuite) Test_DeleteById_WithMissingId() {
	assert.Nil(p.T(), p.repository.DeleteById(uuid.New()))
}

func (p *PaymentRepositoryTestSuite) Test_Update() {
	payment := p.createPayment()

	updatedIncome := model.Payment{
		Id:          payment.Id,
		Name:        fmt.Sprintf("%s-new", payment.Name),
		Description: fmt.Sprintf("%s-new", payment.Description),
		Date:        mocks.Date,
		Sum:         payment.Sum + 100.0,
		HouseId:     payment.HouseId,
		House:       payment.House,
		User:        payment.User,
		UserId:      payment.UserId,
	}

	err := p.repository.Update(updatedIncome)

	assert.Nil(p.T(), err)

	response, err := p.repository.FindById(payment.Id)
	assert.Nil(p.T(), err)
	assert.Equal(p.T(), model.Payment{
		Id:          payment.Id,
		Name:        "Test Payment-new",
		Description: "Test Payment Description-new",
		Date:        updatedIncome.Date,
		Sum:         1100.0,
		HouseId:     payment.HouseId,
		House:       payment.House,
		User:        payment.User,
		UserId:      payment.UserId,
	}, response)
}

func (p *PaymentRepositoryTestSuite) Test_Update_WithMissingId() {
	assert.Nil(p.T(), p.repository.Update(model.Payment{Id: uuid.New()}))
}

func (p *PaymentRepositoryTestSuite) createPayment() model.Payment {
	payment := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id)

	p.CreateEntity(payment)

	return payment
}
