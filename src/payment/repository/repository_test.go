package repository

import (
	"fmt"
	dependencyMocks "github.com/VlasovArtem/hob/src/common/dependency/mocks"
	"github.com/VlasovArtem/hob/src/db"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/payment/mocks"
	"github.com/VlasovArtem/hob/src/payment/model"
	providerMocks "github.com/VlasovArtem/hob/src/provider/mocks"
	providerModel "github.com/VlasovArtem/hob/src/provider/model"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func Test_Initialize(t *testing.T) {
	provider := new(dependencyMocks.DependenciesProvider)
	object := db.Database{}

	provider.On("FindByType", mock.Anything, mock.Anything).Return(&object)

	repository := PaymentRepositoryObject{}

	newObject := repository.Initialize(provider)

	assert.Equal(t, db.modeledDatabase{DatabaseService: &object, Model: model.Payment{}}, newObject.(*PaymentRepositoryObject).database)
}

func Test_GetEntity(t *testing.T) {
	object := PaymentRepositoryObject{}
	assert.Equal(t, model.Payment{}, object.GetEntity())
}

type PaymentRepositoryTestSuite struct {
	database.DBTestSuite
	repository      PaymentRepository
	createdUser     userModel.User
	createdHouse    houseModel.House
	createdProvider providerModel.Provider
}

func (p *PaymentRepositoryTestSuite) SetupSuite() {
	p.InitDBTestSuite()

	p.CreateRepository(
		func(service db.DatabaseService) {
			p.repository = NewPaymentRepository(service)
		},
	).
		AddAfterTest(func(service db.DatabaseService) {
			database.TruncateTable(service, model.Payment{})
		}).
		AddAfterSuite(func(service db.DatabaseService) {
			database.TruncateTable(service, providerModel.Provider{})
			database.TruncateTable(service, houseModel.House{})
			database.TruncateTable(service, userModel.User{})
		}).
		ExecuteMigration(userModel.User{}, houseModel.House{}, providerModel.Provider{}, model.Payment{})

	p.createdUser = userMocks.GenerateUser()
	p.CreateEntity(&p.createdUser)

	p.createdHouse = houseMocks.GenerateHouse(p.createdUser.Id)
	p.CreateEntity(&p.createdHouse)

	p.createdProvider = providerMocks.GenerateProvider(p.createdUser.Id)
	p.CreateEntity(&p.createdProvider)
}

func TestPaymentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentRepositoryTestSuite))
}

func (p *PaymentRepositoryTestSuite) Test_Create() {
	payment := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)

	actual, err := p.repository.Create(payment)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), payment, actual)

	p.Delete(payment)
}

func (p *PaymentRepositoryTestSuite) Test_Creat_WithMissingUser() {
	payment := mocks.GeneratePayment(p.createdHouse.Id, uuid.New(), p.createdProvider.Id)

	actual, err := p.repository.Create(payment)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), payment, actual)
}

func (p *PaymentRepositoryTestSuite) Test_Creat_WithMissingHouse() {
	payment := mocks.GeneratePayment(uuid.New(), p.createdUser.Id, p.createdProvider.Id)

	actual, err := p.repository.Create(payment)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), payment, actual)
}

func (p *PaymentRepositoryTestSuite) Test_Creat_WithMissingProvider() {
	payment := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, uuid.New())

	actual, err := p.repository.Create(payment)

	assert.NotNil(p.T(), err)
	assert.Equal(p.T(), payment, actual)
}

func (p *PaymentRepositoryTestSuite) Test_CreateBatch() {
	first := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	first.Name = "First Income"
	second := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	second.Name = "Second Income"

	actual, err := p.repository.CreateBatch([]model.Payment{first, second})

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), []model.Payment{first, second}, actual)
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
	first := p.createPayment()
	second := p.createPayment()

	actual := p.repository.FindByUserId(p.createdUser.Id, 2, 0, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{second.ToDto(), first.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByUserId_WithFromAndTo() {
	first := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	first.Date = time.Now().AddDate(0, 0, -1)
	p.CreateEntity(&first)

	second := p.createPayment()

	from := time.Now().Add(-time.Hour * 12)
	to := time.Now()
	actual := p.repository.FindByUserId(p.createdUser.Id, 2, 0, &from, &to)

	assert.Equal(p.T(), []model.PaymentDto{second.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByUserId_WithFrom() {
	first := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	first.Date = time.Now().AddDate(0, 0, -1)
	p.CreateEntity(&first)

	second := p.createPayment()

	from := time.Now().Add(-time.Hour * 12)
	actual := p.repository.FindByUserId(p.createdUser.Id, 2, 0, &from, nil)

	assert.Equal(p.T(), []model.PaymentDto{second.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByUserId_WithLimit() {
	_ = p.createPayment()
	second := p.createPayment()

	actual := p.repository.FindByUserId(p.createdUser.Id, 1, 0, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{second.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByUserId_WithOffset() {
	first := p.createPayment()
	_ = p.createPayment()

	actual := p.repository.FindByUserId(p.createdUser.Id, 1, 1, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{first.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByUserId_WithMissingUserId() {
	actual := p.repository.FindByUserId(uuid.New(), 0, 1, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByHouseId() {
	first := p.createPayment()
	second := p.createPayment()

	actual := p.repository.FindByHouseId(p.createdHouse.Id, 2, 0, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{second.ToDto(), first.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByHouseId_WithFromAndTo() {
	first := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	first.Date = time.Now().AddDate(0, 0, -1)
	p.CreateEntity(&first)

	second := p.createPayment()

	from := time.Now().Add(-time.Hour * 12)
	to := time.Now()

	actual := p.repository.FindByHouseId(p.createdHouse.Id, 2, 0, &from, &to)

	assert.Equal(p.T(), []model.PaymentDto{second.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByHouseId_WithFrom() {
	first := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	first.Date = time.Now().AddDate(0, 0, -1)
	p.CreateEntity(&first)

	second := p.createPayment()

	from := time.Now().Add(-time.Hour * 12)

	actual := p.repository.FindByHouseId(p.createdHouse.Id, 2, 0, &from, nil)

	assert.Equal(p.T(), []model.PaymentDto{second.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByHouseId_WithLimit() {
	_ = p.createPayment()
	second := p.createPayment()

	actual := p.repository.FindByHouseId(p.createdHouse.Id, 1, 0, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{second.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByHouseId_WithOffset() {
	first := p.createPayment()
	_ = p.createPayment()

	actual := p.repository.FindByHouseId(p.createdHouse.Id, 1, 1, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{first.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByHouseId_WithMissingId() {
	actual := p.repository.FindByHouseId(uuid.New(), 0, 10, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByProviderId() {
	first := p.createPayment()
	second := p.createPayment()

	actual := p.repository.FindByProviderId(p.createdProvider.Id, 2, 0, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{second.ToDto(), first.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByProviderId_WithFromAndTo() {
	first := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	first.Date = time.Now().AddDate(0, 0, -1)
	p.CreateEntity(&first)

	second := p.createPayment()

	from := time.Now().Add(-time.Hour * 12)
	to := time.Now()

	actual := p.repository.FindByProviderId(p.createdProvider.Id, 2, 0, &from, &to)

	assert.Equal(p.T(), []model.PaymentDto{second.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByProviderId_WithFrom() {
	first := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	first.Date = time.Now().AddDate(0, 0, -1)
	p.CreateEntity(&first)

	second := p.createPayment()

	from := time.Now().Add(-time.Hour * 12)

	actual := p.repository.FindByProviderId(p.createdProvider.Id, 2, 0, &from, nil)

	assert.Equal(p.T(), []model.PaymentDto{second.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByProviderId_WithLimit() {
	_ = p.createPayment()
	second := p.createPayment()

	actual := p.repository.FindByProviderId(p.createdProvider.Id, 1, 0, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{second.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByProviderId_WithOffset() {
	first := p.createPayment()
	_ = p.createPayment()

	actual := p.repository.FindByProviderId(p.createdProvider.Id, 1, 1, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{first.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByProviderId_WithMissingId() {
	actual := p.repository.FindByProviderId(uuid.New(), 0, 1, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{}, actual)
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
	p.createdProvider = providerMocks.GenerateProvider(p.createdUser.Id)
	p.CreateEntity(&p.createdProvider)

	payment := p.createPayment()

	updatedIncome := model.Payment{
		Id:          payment.Id,
		Name:        fmt.Sprintf("%s-new", payment.Name),
		Description: fmt.Sprintf("%s-new", payment.Description),
		Date:        mocks.Date,
		Sum:         payment.Sum + 100.0,
		HouseId:     payment.HouseId,
		UserId:      payment.UserId,
		ProviderId:  payment.ProviderId,
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
		ProviderId:  payment.ProviderId,
		Provider:    payment.Provider,
	}, response)
}

func (p *PaymentRepositoryTestSuite) Test_Update_WithMissingId() {
	assert.Nil(p.T(), p.repository.Update(model.Payment{Id: uuid.New()}))
}

func (p *PaymentRepositoryTestSuite) Test_Delete() {
	payment := p.createPayment()

	assert.Nil(p.T(), p.repository.Delete(payment.Id))
}

func (p *PaymentRepositoryTestSuite) Test_Delete_WithMissingEntity() {
	assert.Nil(p.T(), p.repository.Delete(uuid.New()))
}

func (p *PaymentRepositoryTestSuite) createPayment() model.Payment {
	payment := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	payment.Date = time.Now().Truncate(time.Microsecond)

	p.CreateEntity(payment)

	return payment
}
