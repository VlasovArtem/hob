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
	"log"
	"testing"
)

type PaymentRepositorySchedulerTestSuite struct {
	suite.Suite
	database     db.DatabaseService
	repository   PaymentSchedulerRepository
	createdUser  userModel.User
	createdHouse houseModel.House
}

func (p *PaymentRepositorySchedulerTestSuite) SetupSuite() {
	config := db.NewDefaultDatabaseConfiguration()
	config.DBName = "hob_test"
	p.database = db.NewDatabaseService(config)
	p.repository = NewPaymentSchedulerRepository(p.database)
	err := p.database.D().AutoMigrate(model.PaymentScheduler{})

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

func (p *PaymentRepositorySchedulerTestSuite) TearDownSuite() {
	database.DropTable(p.database.D(), houseModel.House{})
	database.DropTable(p.database.D(), userModel.User{})
	database.DropTable(p.database.D(), model.PaymentScheduler{})
}

func TestPaymentRepositorySchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentRepositorySchedulerTestSuite))
}

func (p *PaymentRepositorySchedulerTestSuite) Test_Create() {
	paymentScheduler := mocks.GeneratePaymentScheduler(p.createdHouse.Id, p.createdUser.Id)

	actual, err := p.repository.Create(paymentScheduler)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), paymentScheduler, actual)
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

	create, err := p.repository.Create(payment)

	assert.Nil(p.T(), err)

	return create
}
