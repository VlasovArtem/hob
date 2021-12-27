package repository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/db"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/payment/scheduler/mocks"
	"github.com/VlasovArtem/hob/src/payment/scheduler/model"
	providerMocks "github.com/VlasovArtem/hob/src/provider/mocks"
	providerModel "github.com/VlasovArtem/hob/src/provider/model"
	"github.com/VlasovArtem/hob/src/scheduler"
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
	repository      PaymentSchedulerRepository
	createdUser     userModel.User
	createdHouse    houseModel.House
	createdProvider providerModel.Provider
}

func (p *PaymentRepositorySchedulerTestSuite) SetupSuite() {
	p.InitDBTestSuite()

	p.CreateRepository(
		func(service db.DatabaseService) {
			p.repository = NewPaymentSchedulerRepository(service)
		},
	).
		AddMigrators(userModel.User{}, houseModel.House{}, providerModel.Provider{}, model.PaymentScheduler{})

	p.createdUser = userMocks.GenerateUser()
	p.CreateConstantEntity(&p.createdUser)

	p.createdHouse = houseMocks.GenerateHouse(p.createdUser.Id)
	p.CreateConstantEntity(&p.createdHouse)

	p.createdProvider = providerMocks.GenerateProvider(p.createdUser.Id)
	p.CreateConstantEntity(&p.createdProvider)
}

func TestPaymentRepositorySchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentRepositorySchedulerTestSuite))
}

func (p *PaymentRepositorySchedulerTestSuite) Test_Create() {
	paymentScheduler := mocks.GeneratePaymentScheduler(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)

	actual, err := p.repository.Create(paymentScheduler)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), paymentScheduler, actual)

	p.Delete(paymentScheduler)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_Creat_WithMissingUser() {
	paymentScheduler := mocks.GeneratePaymentScheduler(p.createdHouse.Id, uuid.New(), p.createdProvider.Id)

	actual, err := p.repository.Create(paymentScheduler)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), paymentScheduler, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_Creat_WithMissingHouse() {
	paymentScheduler := mocks.GeneratePaymentScheduler(uuid.New(), p.createdUser.Id, p.createdProvider.Id)

	actual, err := p.repository.Create(paymentScheduler)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), paymentScheduler, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_Creat_WithMissingProvider() {
	paymentScheduler := mocks.GeneratePaymentScheduler(p.createdHouse.Id, p.createdUser.Id, uuid.New())

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

	assert.Equal(p.T(), []model.PaymentSchedulerDto{payment.ToDto()}, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_FindByUserId_WithMissingId() {
	actual := p.repository.FindByUserId(uuid.New())

	assert.Equal(p.T(), []model.PaymentSchedulerDto{}, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_FindByHouseId() {
	payment := p.createPaymentScheduler()

	actual := p.repository.FindByHouseId(payment.HouseId)

	assert.Equal(p.T(), []model.PaymentSchedulerDto{payment.ToDto()}, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_FindByHouseId_WithMissingId() {
	actual := p.repository.FindByHouseId(uuid.New())

	assert.Equal(p.T(), []model.PaymentSchedulerDto{}, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_FindByProviderId() {
	payment := p.createPaymentScheduler()

	actual := p.repository.FindByProviderId(payment.ProviderId)

	assert.Equal(p.T(), []model.PaymentSchedulerDto{payment.ToDto()}, actual)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_FindByProviderId_WithMissingId() {
	actual := p.repository.FindByProviderId(uuid.New())

	assert.Equal(p.T(), []model.PaymentSchedulerDto{}, actual)
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

func (p *PaymentRepositorySchedulerTestSuite) Test_Update() {
	newProvider := providerMocks.GenerateProvider(p.createdUser.Id)
	p.CreateEntity(&newProvider)

	payment := p.createPaymentScheduler()

	updatedIncome := model.PaymentScheduler{
		Id:          payment.Id,
		Name:        fmt.Sprintf("%s-new", payment.Name),
		Description: fmt.Sprintf("%s-new", payment.Description),
		Sum:         payment.Sum + 100.0,
		Spec:        scheduler.WEEKLY,
		ProviderId:  newProvider.Id,
		HouseId:     payment.HouseId,
		House:       payment.House,
		User:        payment.User,
		UserId:      payment.UserId,
	}

	err := p.repository.Update(updatedIncome)

	assert.Nil(p.T(), err)

	response, err := p.repository.FindById(payment.Id)
	assert.Nil(p.T(), err)
	assert.Equal(p.T(), model.PaymentScheduler{
		Id:          payment.Id,
		Name:        "Test Payment-new",
		Description: "Test Payment Description-new",
		Sum:         1100.0,
		Spec:        scheduler.WEEKLY,
		ProviderId:  newProvider.Id,
		HouseId:     payment.HouseId,
		House:       payment.House,
		User:        payment.User,
		UserId:      payment.UserId,
	}, response)
}

func (p *PaymentRepositorySchedulerTestSuite) Test_Update_WithMissingId() {
	assert.Nil(p.T(), p.repository.Update(model.PaymentScheduler{Id: uuid.New()}))
}

func (p *PaymentRepositorySchedulerTestSuite) createPaymentScheduler() model.PaymentScheduler {
	payment := mocks.GeneratePaymentScheduler(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)

	p.CreateEntity(payment)

	return payment
}
