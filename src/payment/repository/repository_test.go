package repository

import (
	"fmt"
	"github.com/VlasovArtem/hob/src/db"
	groupMocks "github.com/VlasovArtem/hob/src/group/mocks"
	groupModel "github.com/VlasovArtem/hob/src/group/model"
	houseMocks "github.com/VlasovArtem/hob/src/house/mocks"
	houseModel "github.com/VlasovArtem/hob/src/house/model"
	"github.com/VlasovArtem/hob/src/payment/mocks"
	"github.com/VlasovArtem/hob/src/payment/model"
	providerMocks "github.com/VlasovArtem/hob/src/provider/mocks"
	providerModel "github.com/VlasovArtem/hob/src/provider/model"
	"github.com/VlasovArtem/hob/src/test/testhelper"
	"github.com/VlasovArtem/hob/src/test/testhelper/database"
	userMocks "github.com/VlasovArtem/hob/src/user/mocks"
	userModel "github.com/VlasovArtem/hob/src/user/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type PaymentRepositoryTestSuite struct {
	database.DBTestSuite[model.Payment]
	repository      PaymentRepository
	createdUser     userModel.User
	createdHouse    houseModel.House
	createdProvider providerModel.Provider
}

func (p *PaymentRepositoryTestSuite) SetupSuite() {
	p.InitDBTestSuite()

	p.CreateRepository(
		func(service db.ModeledDatabase[model.Payment]) {
			p.repository = NewPaymentRepository(service)
		},
	).
		AddAfterTest(func(service db.ModeledDatabase[model.Payment]) {
			service.DB().Exec("DELETE FROM house_groups")
			testhelper.TruncateTable(service, groupModel.Group{})
			testhelper.TruncateTable(service, model.Payment{})
		}).
		AddAfterSuite(func(service db.ModeledDatabase[model.Payment]) {
			testhelper.TruncateTable(service, providerModel.Provider{})
			testhelper.TruncateTable(service, houseModel.House{})
			testhelper.TruncateTable(service, userModel.User{})
		}).
		ExecuteMigration(userModel.User{}, houseModel.House{}, providerModel.Provider{}, model.Payment{}, groupModel.Group{})

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

func (p *PaymentRepositoryTestSuite) Test_CreateBatch() {
	first := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	first.Name = "First Income"
	second := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	second.Name = "Second Income"

	payments := []model.Payment{first, second}
	err := p.repository.Create(&payments)

	assert.Nil(p.T(), err)
	assert.Equal(p.T(), []model.Payment{first, second}, payments)
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

func (p *PaymentRepositoryTestSuite) Test_FindByGroupId() {
	group := groupMocks.GenerateGroup(p.createdUser.Id)

	firstHouse := houseMocks.GenerateHouse(p.createdUser.Id)
	firstHouse.Groups = []groupModel.Group{group}
	p.CreateEntity(&firstHouse)
	secondHouse := houseMocks.GenerateHouse(p.createdUser.Id)
	secondHouse.Groups = []groupModel.Group{group}
	p.CreateEntity(&secondHouse, "Groups.*")

	firstPayment := mocks.GeneratePayment(firstHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	firstPayment.Date = time.Now().Truncate(time.Microsecond)

	p.CreateEntity(firstPayment)

	secondPayment := mocks.GeneratePayment(secondHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	secondPayment.Date = time.Now().Truncate(time.Microsecond)

	p.CreateEntity(secondPayment)

	actual := p.repository.FindByGroupId(group.Id, -1, -1, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{secondPayment.ToDto(), firstPayment.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByGroupId_WithFromAndTo() {
	group := groupMocks.GenerateGroup(p.createdUser.Id)

	firstHouse := houseMocks.GenerateHouse(p.createdUser.Id)
	firstHouse.Groups = []groupModel.Group{group}
	p.CreateEntity(&firstHouse)
	secondHouse := houseMocks.GenerateHouse(p.createdUser.Id)
	secondHouse.Groups = []groupModel.Group{group}
	p.CreateEntity(&secondHouse, "Groups.*")

	firstPayment := mocks.GeneratePayment(firstHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	firstPayment.Date = time.Now().AddDate(0, 0, -1)

	p.CreateEntity(firstPayment)

	secondPayment := mocks.GeneratePayment(secondHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	secondPayment.Date = time.Now().Truncate(time.Microsecond)

	p.CreateEntity(secondPayment)

	from := time.Now().Add(-time.Hour * 12)
	to := time.Now()

	actual := p.repository.FindByGroupId(group.Id, 2, 0, &from, &to)

	assert.Equal(p.T(), []model.PaymentDto{secondPayment.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByGroupId_WithFrom() {
	group := groupMocks.GenerateGroup(p.createdUser.Id)

	firstHouse := houseMocks.GenerateHouse(p.createdUser.Id)
	firstHouse.Groups = []groupModel.Group{group}
	p.CreateEntity(&firstHouse)
	secondHouse := houseMocks.GenerateHouse(p.createdUser.Id)
	secondHouse.Groups = []groupModel.Group{group}
	p.CreateEntity(&secondHouse, "Groups.*")

	firstPayment := mocks.GeneratePayment(firstHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	firstPayment.Date = time.Now().AddDate(0, 0, -1)

	p.CreateEntity(firstPayment)

	secondPayment := mocks.GeneratePayment(secondHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	secondPayment.Date = time.Now().Truncate(time.Microsecond)

	p.CreateEntity(secondPayment)

	from := time.Now().Add(-time.Hour * 12)

	actual := p.repository.FindByGroupId(group.Id, 2, 0, &from, nil)

	assert.Equal(p.T(), []model.PaymentDto{secondPayment.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByGroupId_WithLimit() {
	group := groupMocks.GenerateGroup(p.createdUser.Id)

	firstHouse := houseMocks.GenerateHouse(p.createdUser.Id)
	firstHouse.Groups = []groupModel.Group{group}
	p.CreateEntity(&firstHouse)
	secondHouse := houseMocks.GenerateHouse(p.createdUser.Id)
	secondHouse.Groups = []groupModel.Group{group}
	p.CreateEntity(&secondHouse, "Groups.*")

	firstPayment := mocks.GeneratePayment(firstHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	firstPayment.Date = time.Now().Truncate(time.Microsecond)

	p.CreateEntity(firstPayment)

	secondPayment := mocks.GeneratePayment(secondHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	secondPayment.Date = time.Now().Truncate(time.Microsecond)

	p.CreateEntity(secondPayment)

	actual := p.repository.FindByGroupId(group.Id, 1, 0, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{secondPayment.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByGroupId_WithOffset() {
	group := groupMocks.GenerateGroup(p.createdUser.Id)

	firstHouse := houseMocks.GenerateHouse(p.createdUser.Id)
	firstHouse.Groups = []groupModel.Group{group}
	p.CreateEntity(&firstHouse)
	secondHouse := houseMocks.GenerateHouse(p.createdUser.Id)
	secondHouse.Groups = []groupModel.Group{group}
	p.CreateEntity(&secondHouse, "Groups.*")

	firstPayment := mocks.GeneratePayment(firstHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	firstPayment.Date = time.Now().Truncate(time.Microsecond)

	p.CreateEntity(firstPayment)

	secondPayment := mocks.GeneratePayment(secondHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	secondPayment.Date = time.Now().Truncate(time.Microsecond)

	p.CreateEntity(secondPayment)

	actual := p.repository.FindByGroupId(group.Id, 1, 1, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{firstPayment.ToDto()}, actual)
}

func (p *PaymentRepositoryTestSuite) Test_FindByGroupId_WithMissingId() {
	actual := p.repository.FindByGroupId(uuid.New(), 0, 1, nil, nil)

	assert.Equal(p.T(), []model.PaymentDto{}, actual)
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

	err := p.repository.Update(payment.Id, updatedIncome)

	assert.Nil(p.T(), err)

	response, err := p.repository.Find(payment.Id)
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

func (p *PaymentRepositoryTestSuite) createPayment() model.Payment {
	payment := mocks.GeneratePayment(p.createdHouse.Id, p.createdUser.Id, p.createdProvider.Id)
	payment.Date = time.Now().Truncate(time.Microsecond)

	p.CreateEntity(payment)

	return payment
}
