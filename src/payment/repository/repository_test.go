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
	"payment/mocks"
	"payment/model"
	"test/testhelper/database"
	"testing"
	userMocks "user/mocks"
	userModel "user/model"
)

type PaymentRepositoryTestSuite struct {
	suite.Suite
	database     db.DatabaseService
	repository   PaymentRepository
	createdUser  userModel.User
	createdHouse houseModel.House
}

func (p *PaymentRepositoryTestSuite) SetupSuite() {
	config := db.NewDefaultDatabaseConfiguration()
	config.DBName = "hob_test"
	p.database = db.NewDatabaseService(config)
	p.repository = NewPaymentRepository(p.database)
	err := p.database.D().AutoMigrate(model.Payment{})

	if err != nil {
		log.Fatal(err)
	}

	p.createdUser = userMocks.GenerateUser()
	err = p.database.Create(&p.createdUser)

	if err != nil {
		log.Fatal(err)
	}

	p.createdHouse = houseMocks.GenerateHouse(p.createdUser.Id)

	err = p.database.Create(&p.createdHouse)

	if err != nil {
		log.Fatal(err)
	}
}

func (p *PaymentRepositoryTestSuite) TearDownSuite() {
	database.DropTable(p.database.D(), houseModel.House{})
	database.DropTable(p.database.D(), userModel.User{})
	database.DropTable(p.database.D(), model.Payment{})
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
